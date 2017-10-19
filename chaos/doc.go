/*
Package failure manages all the failure types. A failure among other properties it is made
from multiple attacks. these attacks will be started for a failure defined duration, after
this defined period the failure will be reverted.

Failure, attacks & injection:

A failure can't be executed by itself, it has all the required information to apply the attacks
but not the environment. A failure that can be executed is named and "Injection", the injection
has the running attacks, locks, timers...
This could be wrapped in:

[Injection]
		└──────[Failure]
				   ├─────[Definition]
				   └─────[Attacks]
				   			 ├─────[Attack1]
							 ├─────...
							 └─────[AttackN]

We have to take into account that there are different kind of attacks, on kind of attack is
a permanen attack based on duration, for example a memory allocation: 'allocate 100MiB of
memory for 60m then free these 100MiB of memory'. On the other hand we have attacks that are
periodic but adding a new failure on top of previous one , for example a memory leak, a memory
leak is a periodic memory allocation attack: allocate 10MiB every 30s for 60m. A
third kind of attack can be a periodic attack not based on the previous state, for example a
CPU spike every of 1m every 20m for 4h.

there are a lot of attacks with different nature and a failure can have multiple attacks at
the same time, taking this into account this will the responsability of the attack itself
(periods, flushing run/state, adding on top of previous run/state of the periodic loop, etc)

configuraiton/defition examples can be check on failure/testdata path

Experiment:

On the other part we have the name Experiment, an experiment is only a group of failures that
have wil be (or are being) applied on targets. Every target will have the same failure
(injection when being applied), the Experiment will match these failuresusing selectors, the
selectors will match the tags of the targets.

This would be:

[Experiment0]
	 ├──[Injection0 (Node0)]
	 │		└──────[Failure0]
	 │				   ├─────[Definition]
	 │				   └─────[Attacks]
	 │				   			 ├─────[Attack0]
	 │							 ├─────...
	 │							 └─────[AttackN]
  	 │
	 └──[Injection0  (Node1)]
			└──────[Failure0]
					   ├─────[Definition]
					   └─────[Attacks]
					   			 ├─────[Attack0]
								 ├─────...
								 └─────[AttackN]


*/
package failure // import "github.com/slok/ragnarok/chaos"
