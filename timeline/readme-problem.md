# Linear Stream Schedule Timeline

Here's a few things to keep in mind before tackling this:

  - Try to keep the time you work on this under two hours. While putting in extra effort is laudable, when evaluating these, we'll attempt to judge progress within the first two hours as a baseline.
  - Try to produce a fully documented and working example before focusing too much on algorithms or minutae.

## Objective

Imagine Twitch channels are just shows on TV. One broadcaster ends their broadcast, another takes over. Life isn't that simple in the livestreaming world though. We want you to make it that simple.

Given the data below, we want you to implement an API with a single endpoint: `/timeline`. Calling this endpoint should produce an array of objects that include a start time, an end time and at least the `channel` from our input, at minimum.

Furthermore, the start and end times should not overlap, but also minimize likeliness of cutting off a broadcast too early if the viewer would switch to the next item in the list. Broadcasts can appear multiple times in the output, but going back into an ongoing broadcast that you have already left should be treated as priority 0, only to occur if there is no other alternative, unless the priority on that broadcast is 10.

Since it is inevitable that at some point the overlap might not mathematically favor the previous broadcaster, we have `priority` - an integer from 0-10 that represents the importance of the broadcast to the caller. The higher this number, the more the algorithm should favor this broadcaster when deciding whether to cut over early, or late.

For example, a broadcast with a priority of 0 should never interrupt a broadcast with a priority of 10. A broadcast with a priority of 4 should minimize the cut off period when the next broadcast is only a 6.

**Important:** some of this logic can be offloaded to the (imaginary) client. If you decide to implement ways in which the client can make these choices, you should explain them in documentation, on a comment in the implementation for example, or via a README.

## Data

Notes on this data:

  - Some of these are scheduled in advance - their `type` is `SCHEDULED`. Scheduled broadcasts have been manually entered by the broadcaster and are not derived from past data. These can be treated as more accurate.

```json
[{
  "id": 1,
  "channel": "twitch.tv/ninja",
  "type": "GUESSTIMATE",
  "startsAt": "Thu Jan 17 10:12:00 PST 2019",
  "endsAt": "Thu Jan 17 16:21:00 PST 2019",
  "priority": 10
}, {
  "id": 2,
  "channel": "twitch.tv/lirik",
  "type": "GUESSTIMATE",
  "startsAt": "Thu Jan 17 08:02:00 PST 2019",
  "endsAt": "Thu Jan 17 13:12:00 PST 2019",
  "priority": 7
}, {
  "id": 3,
  "channel": "twitch.tv/ninja",
  "type": "SCHEDULED",
  "startsAt": "Thu Jan 17 19:00:00 PST 2019",
  "endsAt": "Thu Jan 17 19:30:00 PST 2019",
  "priority": 10
}, {
  "id": 4,
  "channel": "twitch.tv/summit1g",
  "type": "GUESSTIMATE",
  "startsAt": "Thu Jan 17 13:00:00 PST 2019",
  "endsAt": "Thu Jan 17 15:00:00 PST 2019",
  "priority": 2 
}, {
  "id": 5,
  "channel": "twitch.tv/annemunition",
  "type": "GUESSTIMATE",
  "startsAt": "Thu Jan 17 04:49:00 PST 2019",
  "endsAt": "Thu Jan 17 11:33:00 PST 2019",
  "priority": 8
}, {
  "id": 6,
  "channel": "twitch.tv/annemunition",
  "type": "GUESSTIMATE",
  "startsAt": "Thu Jan 17 13:01:00 PST 2019",
  "endsAt": "Thu Jan 17 18:13:00 PST 2019",
  "priority": 8
}, {
  "id": 7,
  "channel": "twitch.tv/shroud",
  "type": "GUESSTIMATE",
  "startsAt": "Thu Jan 17 18:52:00 PST 2019",
  "endsAt": "Thu Jan 17 22:41:00 PST 2019",
  "priority": 0
}, {
  "id": 8,
  "channel": "twitch.tv/admiralbahroo",
  "type": "SCHEDULED",
  "startsAt": "Thu Jan 17 7:00:00 PST 2019",
  "endsAt": "Thu Jan 17 14:00:00 PST 2019",
  "priority": 4
}, {
  "id": 9,
  "channel": "twitch.tv/fyrn",
  "type": "GUESSTIMATE",
  "startsAt": "Thu Jan 17 02:30:00 PST 2019",
  "endsAt": "Thu Jan 17 03:30:00 PST 2019",
  "priority": 0
}]
```