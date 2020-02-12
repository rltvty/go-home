# Linear Stream Schedule Timeline

This is a solution to the Linear Stream Schedule Timeline problem, presented to me as an interview challenge.

For more details on the problem, see: [Problem](./readme-problem.md)

## Eric Pinzur - Solution

My general approach is:
1. parse the input data into a strongly-typed golang struct  
1. figure out the start and end of the schedule
1. iterate through the schedule on a per-minute basis
1. for each minute:
   1. get the potential stream choices (ideally excluding previously watched streams), sorted by priorty and stream type
   1. try to determine it should stay on the current stream, or move to the next one based on:
      * if no current stream, then start best potential
      * if current stream priority is 10, stay on it
      * if new stream priority is 10, move to it
      * otherwise, move to new stream is priority is more than 4 higher than current stream

I wrote my solution to the problem in go-lang.  

Unfortunatley, the solution doesn't currently produce the ideal results... at least not compared to what I think they should be if I was doing this with pen and paper.  But I've spent more than 3 hours on this already, and I'm going to stop at this point.

Note that I'm still very much a go-lang beginner, which added to the time it took me to get to this point. 

### Running

switch to the `timeline` folder, get dependcies, and run:

```
cd timeline
go get .
go run .
```

### Code

[`./schedule/schedule.go`](./schedule/schedule.go): contains the code for the `schedule` package.  All the interesting logic is in this file.

[`./schedule/schedule_suite_test.go`](./schedule/schedule_suite_test.go): sets up the ginkgo tests

[`./schedule/schedule_test.go`](./schedule/schedule_test.go): contains the bare minimum of tests for of the sub methods used in the `schedule` package.

[`./main.go`](./main.go): creates the http router and builds the timeline endpoint, using the `schedule` package

### Results

These are the current results I get when doing a GET request on `http://localhost:4123/timeline`:

```
[ 
    { 
       "channel":"twitch.tv/fyrn",
       "streamId":9,
       "startsAt":"2019-01-17T02:30:00Z",
       "endsAt":"2019-01-17T03:30:00Z"
    },
    { 
       "channel":"twitch.tv/annemunition",
       "streamId":5,
       "startsAt":"2019-01-17T04:49:00Z",
       "endsAt":"2019-01-17T11:33:00Z"
    },
    { 
       "channel":"twitch.tv/ninja",
       "streamId":1,
       "startsAt":"2019-01-17T11:34:00Z",
       "endsAt":"2019-01-17T16:21:00Z"
    },
    { 
       "channel":"twitch.tv/annemunition",
       "streamId":6,
       "startsAt":"2019-01-17T16:22:00Z",
       "endsAt":"2019-01-17T18:13:00Z"
    },
    { 
       "channel":"twitch.tv/shroud",
       "streamId":7,
       "startsAt":"2019-01-17T18:52:00Z",
       "endsAt":"2019-01-17T19:00:00Z"
    },
    { 
       "channel":"twitch.tv/ninja",
       "streamId":3,
       "startsAt":"2019-01-17T19:00:00Z",
       "endsAt":"2019-01-17T19:30:00Z"
    },
    { 
       "channel":"twitch.tv/shroud",
       "streamId":7,
       "startsAt":"2019-01-17T19:31:00Z",
       "endsAt":"2019-01-17T22:41:00Z"
    }
 ]
 ```

