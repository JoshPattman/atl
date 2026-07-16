# ATL - A toy protocol for transmitting data over audio
ATL (Audio Transport Language) is a simple protocol that lets you convert any data into audio. It works by allocating 4 audio frequencies per channel. Then, it sends a sync message by playing all 4 frequencies at once. Finally, it sends the data by alternating pairs of frequencies, where one in each pair represents low and on high. Frequency pairs are alternated every bit. This makes ATL self-clocking and a bit more robust to time variations and noise, but quite slow.

To try it out, have a look at `cmd/atlgen`