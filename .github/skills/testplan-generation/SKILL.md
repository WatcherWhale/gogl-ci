---
name: testplan-generation
description: Generate gogl-ci TestPlan files (YAML + Go test scripts) for validating GitLab CI pipelines. Use this skill when asked to create, write, or generate a test plan, validation spec, or test script for a .gitlab-ci.yml pipeline.
---

## What you are generating

A gogl-ci **TestPlan** is a YAML file (plus optional Go script files) that validates a `.gitlab-ci.yml` pipeline. TestPlans are run with `gogl test <pipeline-file> <testplan-file>`.

There are two API versions. Prefer **v1alpha2** for new plans — it is more expressive. Use **v1alpha1** only for simple presence/dependency checks when the user explicitly prefers it.

---

## v1alpha2 TestPlan

### YAML structure

```yaml
apiVersion: test.gogl.ci/v1alpha2
kind: TestPlan
metadata:
  name: <human-readable plan name>

# validations: declarative checks — pipeline is valid (all needs satisfied) for a given context
validations:
  - name: <scenario name>
    defaultBranch: "main"       # sets CI_DEFAULT_BRANCH
    branch: "main"              # sets CI_COMMIT_BRANCH (omit for tag/MR pipelines)
    tag: ""                     # sets CI_COMMIT_TAG (omit when not a tag pipeline)
    mr: false                   # true sets CI_PIPELINE_SOURCE=merge_request_event
    variables:                  # any extra CI variables to inject
      MY_VAR: "value"

# tests: Go scripts interpreted at runtime; each script file is relative to this YAML file
tests:
  - name: <human-readable test group name>
    test: <relative-path-to-go-script>.go
```

**Rules for `validations`:**
- Each entry is a separate pipeline execution context (branch push, MR, tag).
- It only checks that all active jobs have their `needs:` satisfied — it does **not** check custom logic.
- Omit `branch` when testing a tag pipeline; omit `tag` when testing a branch pipeline.

**Rules for `tests`:**
- The `.go` file path is relative to the directory containing the YAML file.
- Multiple test entries can point to different `.go` files; all are run and results aggregated.

### Go test script structure

Test scripts are **interpreted at runtime via yaegi** — they are not compiled into the binary. Every function whose name starts with `Test` and has the signature `func(gitlab.Pipeline) (bool, string)` is automatically discovered and run.

```go
package <any_package_name>

import (
    "github.com/watcherwhale/gogl-ci/pkg/gitlab"
    "github.com/watcherwhale/gogl-ci/pkg/graph"
)

// TestSomething describes what this test verifies.
// Return (true, "") on success; (false, "reason") on failure.
func TestSomething(pipeline gitlab.Pipeline) (bool, string) {
    vars := gitlab.NewVariables().
        WithDefaultBranch("main").
        WithBranch("main")
        // .WithMergeRequest()   — adds MR-related CI variables
        // .WithBranch("feat")   — sets CI_COMMIT_BRANCH

    g, err := graph.NewGraph(pipeline, vars)
    if err != nil {
        return false, err.Error()
    }

    if err := g.Validate(); err != nil {
        return false, err.Error()
    }

    if !g.HasJob("deploy") {
        return false, "deploy job should be present on main branch"
    }

    if !g.HasDependency("build", "deploy") {
        return false, "deploy should depend on build"
    }

    return true, ""
}
```

Only these packages are available inside test scripts (extracted via yaegi):
- `github.com/watcherwhale/gogl-ci/pkg/gitlab`
- `github.com/watcherwhale/gogl-ci/pkg/graph`
- `github.com/watcherwhale/gogl-ci/pkg/rules/interpreter`

Do **not** import standard library packages or any other third-party packages in test scripts.

### Key API — `gitlab.Variables`

```go
gitlab.NewVariables()                    // empty Variables map
vars.WithDefaultBranch("main")          // CI_DEFAULT_BRANCH
vars.WithBranch("feature/x")            // CI_COMMIT_BRANCH
vars.WithMergeRequest()                 // CI_PIPELINE_SOURCE=merge_request_event + MR vars
// tag pipeline: pass map[string]string{"CI_COMMIT_TAG": "v1.0.0"} directly to graph.NewGraph
```

### Key API — `graph.JobGraph`

```go
g, err := graph.NewGraph(pipeline, vars)  // builds the active-job DAG
g.Validate()                              // error if any needs: reference a missing job
g.HasJob("job-name")                      // true if job is active (when != never)
g.HasDependency("dep-job", "job")         // true if job transitively depends on dep-job
g.GetDependencies("job")                  // []string of all transitive dependency job names
g.GetJob("job-name")                      // (gitlab.Job, error)
```

### Key API — `gitlab.Pipeline`

```go
pipeline.GetJobs()                        // map[string]Job — excludes hidden (.) jobs
pipeline.GetActiveJobs(vars)              // map[string]Job — only when != never
pipeline.GetJobsByStage("stage-name")     // map[string]Job
pipeline.Stages                           // []string
pipeline.Variables                        // map[string]string — pipeline-level variables
```

### Key API — `gitlab.Job`

```go
job.Name        // string
job.Stage       // string
job.When        // "on_success" | "on_failure" | "always" | "never" | "manual" | "delayed"
job.Needs       // gitlab.Needs — job.Needs.NoNeeds == true means no explicit needs (uses stage ordering)
job.Variables   // map[string]string
job.Rules       // []gitlab.Rule
```

---

## v1alpha1 TestPlan (simple presence + dependency checks only)

Use when only checking job presence and `dependsOn` relationships for a single pipeline context.

```yaml
apiVersion: test.gogl.ci/v1alpha1
kind: TestPlan
metadata:
  name: <plan name>
  labels:           # optional
    key: value
spec:
  pipeline:
    defaultBranch: "main"
    branch: "main"      # omit for tag pipelines
    tag: ""             # set for tag pipelines
    mr: false           # true = MR pipeline
    variables:
      KEY: "value"
  tests:
    - name: <test description>
      job: <job-name>
      present: true      # default true; set false to assert job is absent
      dependsOn:         # omit to skip dependency check
        - other-job      # assert job transitively depends on these
```

`dependsOn: []` (explicit empty list) asserts the job has **no** dependencies.

---

## Workflow for generating a test plan

1. **Read the `.gitlab-ci.yml`** to understand stages, jobs, rules, and needs.
2. **Identify the pipeline contexts** to test (e.g., push to main, MR, tag release, feature branch).
3. **For each context, decide what to assert:**
   - Which jobs should be active (present/absent)?
   - Which dependency chains must hold (`HasDependency`)?
   - Is the pipeline structurally valid (`g.Validate()`)?
4. **Choose the format:**
   - Structural validity only → `validations:` in v1alpha2 (or v1alpha1 with just `present: true`)
   - Custom assertions → `tests:` in v1alpha2 with a Go script
5. **Name files predictably:** `<pipeline-name>.testplan.yaml` for the YAML and `<scenario>.go` for scripts, placed alongside each other.

## Common patterns

**Assert a job runs only on main:**
```go
func TestJobOnlyOnMain(pipeline gitlab.Pipeline) (bool, string) {
    mainVars := gitlab.NewVariables().WithDefaultBranch("main").WithBranch("main")
    featVars := gitlab.NewVariables().WithDefaultBranch("main").WithBranch("feature/x")

    mainGraph, _ := graph.NewGraph(pipeline, mainVars)
    featGraph, _ := graph.NewGraph(pipeline, featVars)

    if !mainGraph.HasJob("deploy-prod") {
        return false, "deploy-prod should run on main"
    }
    if featGraph.HasJob("deploy-prod") {
        return false, "deploy-prod should NOT run on feature branches"
    }
    return true, ""
}
```

**Assert a full dependency chain:**
```go
func TestDeployChain(pipeline gitlab.Pipeline) (bool, string) {
    vars := gitlab.NewVariables().WithDefaultBranch("main").WithBranch("main")
    g, err := graph.NewGraph(pipeline, vars)
    if err != nil {
        return false, err.Error()
    }

    chain := []string{"lint", "test", "build", "deploy"}
    for i := 1; i < len(chain); i++ {
        if !g.HasDependency(chain[i-1], chain[i]) {
            return false, chain[i] + " must depend on " + chain[i-1]
        }
    }
    return true, ""
}
```

**Tag release pipeline:**
```go
func TestTagRelease(pipeline gitlab.Pipeline) (bool, string) {
    vars := map[string]string{
        "CI_DEFAULT_BRANCH": "main",
        "CI_COMMIT_TAG":     "v1.2.3",
    }
    g, err := graph.NewGraph(pipeline, vars)
    if err != nil {
        return false, err.Error()
    }
    if !g.HasJob("publish") {
        return false, "publish job must run on tag"
    }
    return true, ""
}
```
