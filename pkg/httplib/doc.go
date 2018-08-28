// Package httplib provides convenient http client and methods.
//
// The main goal of this package is to provide a client which is properly set up.
/*
Example usages:
    // Create a client.
    client := htplib.NewClient()

    // Create a client using extended constructor.
    Client := httplib.NewClientWithSettings(
        30 * time.Second, 30 * time.Second,
        90 * time.Second, 10 * time.Second, 1 * time.Second,
        180 * time.Second,
        true, true,
    )
*/
package httplib
