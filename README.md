# Leader Election in Go
[![Build Status](https://travis-ci.org/utsavgupta/go-leader-election.svg?branch=main)](https://travis-ci.org/utsavgupta/go-leader-election)

This project demonstrates leader election implementation in Go. The application uses Google Datastore for persisting data.

For more details on the design of this demo please refer to the companion article found [here](https://www.utsavgupta.in/blog/leader-election/).

## Getting and building the application

```bash
$ git clone https://github.com/utsavgupta/go-leader-election.git
$ cd go-leader-election
$ go build -race
```

## Running the application

You need to have the Google Cloud SDK installed on your system. Steps required to install the same can be found [here](https://cloud.google.com/sdk/docs/quickstart).

Once you have installed the SDK, run the following commands to install and run the Google Datastore emulator.

```bash
$ gcloud components install cloud-datastore-emulator
$ gcloud beta emulators datastore start  
```

In a new terminal window export environment variables to set application stage and port number. In addition to those variables, you need to set the variables needed by the application to initialize the Datastore client. These variables can be set through a gcloud command.

```bash
$ # in the go-leader-election directory
$ export APP_PORT=3000 # change the port number for running additional instances
$ export APP_STAGE=local
$ $(gcloud beta emulators datastore env-init)
$ ./go-leader-election
```

## Verifying leader election

To verify leader election, run multiple instances of the application from their respective terminal windows. The following is the sample output of two schedulers that were run concurrently.

Terminal 1
```
{"level":"info","message":"Started server in 1ms. Listening to requests on port 3000.","application":"go-leader-election","environment":"local","caller":"C:/Users/gupta/Code/go-leader-election/run.go:67","ts":"1604903149067809800"}
{"level":"info","message":"Starting go_preacher scheduler on node boring_morse5","application":"go-leader-election","environment":"local","caller":"C:/Users/gupta/Code/go-leader-election/schedulers/schedulers.go:22","ts":"1604903149067809800"}
{"level":"info","message":"Starting scala_preacher scheduler on node boring_morse5","application":"go-leader-election","environment":"local","caller":"C:/Users/gupta/Code/go-leader-election/schedulers/schedulers.go:22","ts":"1604903149067809800"}
Go is the best modern programming language
Go is the best modern programming language
Go is the best modern programming language
Go is the best modern programming language
Go is the best modern programming language
{"level":"error","message":"Sitting out: Node boring_morse5 could not become leader for job scala_preacher","application":"go-leader-election","environment":"local","caller":"C:/Users/gupta/Code/go-leader-election/schedulers/schedulers.go:31","ts":"1604903806156463400"}
Go is the best modern programming language
Go is the best modern programming language
Go is the best modern programming language
Go is the best modern programming language
Go is the best modern programming language
{"level":"error","message":"Sitting out: Node boring_morse5 could not become leader for job scala_preacher","application":"go-leader-election","environment":"local","caller":"C:/Users/gupta/Code/go-leader-election/schedulers/schedulers.go:31","ts":"1604903897316993000"}
Go is the best modern programming language
Go is the best modern programming language
Go is the best modern programming language
Go is the best modern programming language
Go is the best modern programming language
```

Terminal 2
```
{"level":"info","message":"Started server in 2ms. Listening to requests on port 3001.","application":"go-leader-election","environment":"local","caller":"C:/Users/gupta/Code/go-leader-election/run.go:67","ts":"1604902957584763200"}
{"level":"info","message":"Starting go_preacher scheduler on node ecstatic_meninsky7","application":"go-leader-election","environment":"local","caller":"C:/Users/gupta/Code/go-leader-election/schedulers/schedulers.go:22","ts":"1604902957584763200"}
{"level":"info","message":"Starting scala_preacher scheduler on node ecstatic_meninsky7","application":"go-leader-election","environment":"local","caller":"C:/Users/gupta/Code/go-leader-election/schedulers/schedulers.go:22","ts":"1604902957584763200"}
Scala is the most exciting language on the JVM
Scala is the most exciting language on the JVM
Scala is the most exciting language on the JVM
Scala is the most exciting language on the JVM
Scala is the most exciting language on the JVM
{"level":"error","message":"Idle: Node ecstatic_meninsky7 could not become leader for job go_preacher","application":"go-leader-election","environment":"local","caller":"C:/Users/gupta/Code/go-leader-election/schedulers/schedulers.go:31","ts":"1604903768373642600"}
Scala is the most exciting language on the JVM
Scala is the most exciting language on the JVM
Scala is the most exciting language on the JVM
Scala is the most exciting language on the JVM
Scala is the most exciting language on the JVM
{"level":"error","message":"Idle: Node ecstatic_meninsky7 could not become leader for job go_preacher","application":"go-leader-election","environment":"local","caller":"C:/Users/gupta/Code/go-leader-election/schedulers/schedulers.go:31","ts":"1604903848573250900"}
Scala is the most exciting language on the JVM
Scala is the most exciting language on the JVM
Scala is the most exciting language on the JVM
Scala is the most exciting language on the JVM
Scala is the most exciting language on the JVM
{"level":"error","message":"Idle: Node ecstatic_meninsky7 could not become leader for job go_preacher","application":"go-leader-election","environment":"local","caller":"C:/Users/gupta/Code/go-leader-election/schedulers/schedulers.go:31","ts":"1604903913185202200"}
Scala is the most exciting language on the JVM
Scala is the most exciting language on the JVM
Scala is the most exciting language on the JVM
Scala is the most exciting language on the JVM
Scala is the most exciting language on the JVM
```