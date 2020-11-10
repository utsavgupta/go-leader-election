# Leader Election in Go

This project demonstrates how to implement leader election in Go. This application uses Google Datastore for persisting data.

For more details on the design of this demo please refer to the companion article found [here](https://www.utsavgupta.in/blog/leader-election-go/).

## Getting and Building the Application

```bash
$ git clone https://github.com/utsavgupta/go-leader-election.git
$ cd go-leader-election
$ go build -race
```

## Running the Application

You need to have the Google Cloud SDK installed on your system. Steps required to install the same can be found [here](https://cloud.google.com/sdk/docs/quickstart).

Once you have installed the SDK, run the following commands to install and run the Google Datastore emulator.

```bash
$ gcloud components install cloud-datastore-emulator
$ gcloud beta emulators datastore start  
```

In a new terminal window set environment variables to specify application stage and port number. In addition to that you need to set the variables needed by the application to initialize the Datastore client. These variables can be set through a gcloud command.

```bash
$ # in the go-leader-election directory
$ export APP_PORT=3000
$ export APP_STAGE=local
$ $(gcloud beta emulators datastore env-init)
$ ./go-leader-election
```

In order to run multiple instances of the application and see leader election work, run the same application on different ports.