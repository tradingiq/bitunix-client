# Balance Channel Subscription

## Description

Balance

## Push Parameters

| Parameter         |                                                                        | Type         | Description                            |
|-------------------|------------------------------------------------------------------------|--------------|----------------------------------------|
| ch                | #broken description shows channel name "position" rather than "balance" | String       | Channel name: position                 |
| ts                |                                                                        | Int64        | Time stamp                             |
| data              |                                                                        | List<Object> | Subscription data                      |
| > coin            |                                                                        | String       | coin                                   |
| > available       |                                                                        | String       | available                              |
| > frozen          |                                                                        | String       | frozen = isolationFrozen + crossFrozen |
| > isolationFrozen |                                                                        | String       | Freeze on a per warehouse basis        |
| > crossFrozen     |                                                                        | String       | Full warehouse entrusted freeze        |
| > margin          |                                                                        | String       | Margin                                 |
| > isolationMargin |                                                                        | String       | Warehouse by warehouse margin          |
| > crossMargin     |                                                                        | String       | Full warehouse margin                  |
| > expMoney        |                                                                        | String       | Experience Money                       |