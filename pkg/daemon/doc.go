// Package daemon provides a convenient way to execute a process.
//
/*
Daemon starts anything which implements Runner interface.
It listens either for an error from the Runner or for a signal from OS.

The package provides to ways to use a Daemon:

1. Using built-in wrapper Service to fulfill the requirements of the Runner:

    // Create a service for daemon to run.
    d := daemon.NewDaemon()
    f := func() error {
        err := doSomething()
        return err
    }
    // Set f() as Cmd for a Service.
    s := daemon.Service{Cmd: f}

    // Set process for the Daemon.
    d.SetProcess(s)

    // Or just
    d.SetProcess(daemon.Service{Cmd: f})

    // Start the Daemon.
    if err := d.Start(); err != nil {
        return err
    }

2. Using a service which fulfills the Runner by itself:

    // Create a service for daemon to run.
    d := daemon.NewDaemon()

    // Create your service.
    a := &agent.Config{...}

    // Agent implements Runner by providing (Run() error) method.
    // We can pass it as Runner to SetProcess.
    d.SetProcess(daemon.Runner(a))

    // Start the Daemon.
    if err := d.Start(); err != nil {
        return err
    }

The recommended way to use this package is:
- create an instance of daemon
- create a children context from daemon's context
- create an instance of your service passing that child context
- inside your service you should listen for the context to be cancelled in order to finish properly
*/
package daemon
