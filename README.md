# GoGl-ci - Getting insight into your GitLab pipelines

GoGl is a tool for parsing and testing CI/CD pipelines. It can (almost) fully parse a GitLab CI pipeline with includes support.

## What usecase it solves

Writing pipelines is often paired with high uncertainty and countless commits.
In some cases you can only really test if your pipeline works correctly when it
has been actually merged into the default branch.

With GoGl you can write TestPlans that define the things you want to be certain
of before merging into a default branch. You can test when jobs will run and if
these jobs have dependencies on others. This will make sure pipelines always
execute correctly.

## Project Status

### Parsing Status

Implemented keywords:

- [x] Default
- [ ] Include
    - [x] Local
    - [x] Project
    - [x] Remote
    - [ ] Component
    - [ ] Template
    - [ ] Inputs
- [x] Stages
- [ ] Workflows
- [ ] Header Keywords
    - [ ] Spec
        - [ ] Inputs
- [ ] Job Keywords
    - [x] allow_failure
        - [x] exit_codes
    - [ ] artifacts
        - [ ] paths
        - [ ] exclude
        - [ ] expire_in
        - [ ] expose_as
        - [ ] name
        - [ ] public
        - [ ] access
        - [ ] reports
        - [ ] untracked
        - [ ] when
    - [ ] cache
        - [ ] paths
        - [ ] key
        - [ ] key:files
        - [ ] key:prefix
        - [ ] untracked
        - [ ] unprotect
        - [ ] when
        - [ ] policy
        - [ ] fallback_keys
    - [x] coverage
    - [ ] dast_configuration
    - [ ] dependencies
    - [ ] environment
        - [ ] name
        - [ ] url
        - [ ] on_stop
        - [ ] action
        - [ ] auto_stop_in
        - [ ] kubernetes
        - [ ] deployment_tier
        - [ ] Dynamic environments
    - [x] extends
    - [ ] hooks
        - [ ] pre_get_sources_script
    - [ ] identity
    - [ ] id_tokens
    - [ ] image
        - [x] name
        - [x] entrypoint
        - [ ] docker
        - [ ] pull_policy
    - [ ] inherit
        - [ ] default
        - [ ] variables
    - [x] interruptible
    - [ ] needs
        - [x] artifacts
        - [x] project
        - [ ] pipeline:job
        - [x] optional
        - [ ] pipeline
        - [ ] parallel:matrix
    - [ ] pages
        - [ ] publish
        - [ ] pages.path_prefix
    - [ ] parallel
        - [ ] parallel:matrix
    - [ ] release
        - [ ] tag_name
        - [ ] tag_message
        - [ ] name
        - [ ] description
        - [ ] ref
        - [ ] milestones
        - [ ] released_at
        - [ ] assets:links
    - [ ] resource_group
    - [ ] retry
        - [ ] when
        - [ ] exit_codes
    - [ ] rules
        - [x] if
        - [ ] changes
        - [ ] changes:paths
        - [ ] changes:compare_to
        - [ ] exists
        - [ ] exists:paths
        - [ ] exists:project
        - [x] when
        - [x] allow_failure
        - [ ] needs
        - [ ] variables
        - [ ] interruptible
    - [x] script, before_script, after_script
        - [ ] Reference Tags
    - [ ] secrets
        - [ ] vault
        - [ ] gcp_secret_manager
        - [ ] azure_key_vault
        - [ ] file
        - [ ] token
    - [ ] services
        - [ ] docker
        - [ ] pull_policy
    - [x] stage
    - [ ] tags
    - [ ] timeout
    - [ ] trigger
        - [ ] include
        - [ ] project
        - [ ] strategy
        - [ ] forward
    - [x] variables
        - [ ] description
        - [ ] value
        - [ ] options
        - [ ] expand
    - [ ] when
- Deprecated keywords, these won't be supported
    - Globally-defined image, services, cache, before_script, after_script
    - only / except
        - only:refs / except:refs
        - only:variables / except:variables


### Rules Interpreter Status

- [x] Variables
- [x] Strings
- [x] null
- [x] Equality Operators
    - [x] Regex
- [x] Logical Operators
- [x] String Null Check
