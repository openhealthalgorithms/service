# Contributing Guideline

## Topics

* [The Flow](#the-flow)
* [Cross-Platform Code](#cross-platform-code)
* [Committing](#committing)
* [External Dependencies](#external-dependencies)
* [Tools](#tools)
* [Coding Style](#coding-style)

## The Flow

We follow The Git Flow.

You can read about it in details [here](https://www.atlassian.com/git/tutorials/comparing-workflows/gitflow-workflow).

### Branching Rules

* features are being developed in a separate branches derived from `develop`
* hotfixes are being developed in a separate branches derived from `master`
* tags are being used when it's suitable. For example, we can tag some commit to fix the state in the history
* **any** merge into `master` and `develop` must be done with a Pull Request
* `develop` is **always** working
* `master` is **always** stable

## Cross-Platform Code

This section describes our approach to develop cross-platform codebase.

### Code

The Service itself is written with full support for the different platforms.  
The general recommendation is to avoid platform-specific code outside the set of installers.

### Files

As you know `golang` has two approaches to maintain a cross-platform code:

* using `build` tag
* using suffixes for filenames

The Project uses both approaches.  
The rules are listed below are **must** to follow.

* by default we do not split the code for platforms
* if the implementation differs between any posix compatible system and windows:

  * `packagename.go` must be provided **with no** build tags
  * `packagename_posix.go` must be provided for other platforms **with** build tag `// +build linux darwin`
  * `packagename_windows.go` must be provided for windows **with no** build tags

* if the implementation differs between mixed platforms:

  * `packagename.go` must be provided **with no** build tags
  * `packagename_pltf1_pltf2.go` must be provided **with** build tag `// +build platform1 platform2`
  * `packagename_platform.go` must be provided **with no** build tags for the different platform

* if the implementation differs between many platforms and architectures:

  * `packagename.go` must be provided **with no** build tags
  * `packagename_platform1.go` must be provided **with no** build tags
  * `packagename_platform2.go` must be provided **with no** build tags
  * `packagename_platform3_arch1.go` must be provided **with no** build tags
  * `packagename_platform3_arch2.go` must be provided **with no** build tags

## Committing

You can read about the topic [here](https://chris.beams.io/posts/git-commit/).

### Committing Rules

* the first word of a commit message must start with a Capital
* the length of the subject line should be near 50 symbols
* use the imperative mood in the subject line
* the period at the end of the subject line is prohibited
* `exception` it is considered `ok` to have more than one commit with the same subject line (repeatable commits within the same task)

The correct subject line of a commit message must complete the following sentence:

> If applied, this commit will _your subject line here_

**The Correct Example**:

> If applied, this commit will _Refactor subsystem X for readability_

**The INCORRECT Example**:

> If applicaed, this commit will _fixed bug with Y_

Or

> If applied, this commit will _changing behavior of Z_

## External Dependencies

We use the very restricted set of dependencies.

Our base approach is to rely on standard library as much as possible.  
It doesn't mean that we do not use them at all.  Instead, we use and follow the rules.  

When to use a dependency:

* when it **simplifies** the solution of the task. Good examples:
  * use [`gopsutil`](https://github.com/shirou/gopsutil) to gather the information about the system
* when it would take a lot of time to implement the reliable soulution for all the cases:
  * use [`govalidator`](https://github.com/asaskevich/govalidator) to match a certain string/numeric values for common formats to avoid creating own matchers

When not to use a dependency:

* when the standard library is enough to solve the task

The requirements for a dependency you're about do add and use:

* the last activity in the repository of the library not far than `Early 2017`. We do not rely on outdated stuff. It just won't compile once

Of course, there are exceptions.

## Tools

We use the restricted set of tools.

The intention is to keep the process as simple as possible.

## Coding Style

The list of rules partially taken from the Docker repository and extended with our set.

### Rules

1. All code should be formatted with `gofmt`.
2. All code should pass the default levels of [`golint`](https://github.com/golang/lint).
3. All code should follow the guidelines covered in [Effective Go](http://golang.org/doc/effective_go.html) and [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments).
    1. This is Golang. So do not try to use well known approaches from other languages. Think in Golang.
4. Comment the code. Tell us what and why.
5. Document _all_ declarations and methods, even private ones. Declare expectations, caveats and anything else that may be important.
6. Variable name length should be proportional to its context and no longer. `thisIsNotJavaOrEvenPythonOrCSharp`. In practice, short methods will have short variable names and globals will have longer names.
7. Counterintuitive, but `tools` package is allowed here. For example, for `UnzipFile` func. Why to repeat it many times? :).
8. All tests should run with `go test` and outside tooling should not be required. No, we don't need another unit testing framework.

### Additional set of rules, to emphasize the information from Guidelines

1. Named return values are limited to use and should be used only in case of `defer`red modification.
2. Using slices by pointers is strictly prohibited unless you're passing it to a method which can modify its length. In such case it is almost a must.
3. Global variables should be avoided in general. But there is certain amount of exceptions. See the code.
4. Global default objects should be avoided in general. For example, it is better not to use `http.DefaultServeMux`.
5. Related types and their methods should be placed in the same file.
6. The most important types and structures should be close to the beginning of the file. The order defines importance.
7. Declare interfaces carefully! We're not using them unless we must represent a set of objects.
8. Beforehand creation of Interfaces is strictly prohibited. An Interface must be created only when it is really necessary.
9. The usage of `log.Fatal` and `log.Fatalf` is strictly prohibited. Why - see the next point.
10. The usage of `os.Exit` is restricted. It must be called only as the first deferred function as it's shown in the `main.go`. Any call to `os.Exit` exists immediately without calling any deferred functions. See [here](https://golang.org/pkg/os/#Exit).
    1. `log.Fatal` and `log.Fatalf` use `os.Exit` internally.
