# State Diagram

            +----------------------------------+
            |                                  |
            |            Root State            |
            |                                  |
            +----------------------------------+
            |                                  |
            |          State Headers:          |
            | {Call Stack, Memory, Current...} |
            |                                  |
            +----------------+-----------------+
                             |
        +------------------------------------------+
        |                    |                     |
        |                    |                     |
        |                    |                     |
    +------+------+      +------+------+       +------+------+
    |             |      |             |       |             |
    | State Child |      | State Child |       | State Child |   <---+
    |             |      |             |       |             |       |
    +-------------+      +-------------+       +------+------+       |
                                                      |              |
           ^                    ^                     |              |
           |                    |                     |              |
           |                    |                     |              |
           |                    |              +------+------+       |
           +                    +              |             |       |
       Contested            Contested          | State Child |   <---+
       Data-Set             Data+Set           |             |       |
        (e.g.             (client still        +------+------+       |
       hardfork)           out of date)               |              |
                                                      |              |
                                                      |              |
                                                      |              |
                                               +------+------+       |
                                               |             |       |
                                               | State Child |   <---+
                                               |             |       |
                                               +-------------+       +
                                                                 Uncontested
                                                                 Data+Set
                                                             (e.g. successful
                                                                 hard+fork)
