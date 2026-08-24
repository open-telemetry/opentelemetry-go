# Experimental Exemplar Reservoirs

This package contains experimental exemplar reservoirs for the OpenTelemetry Go SDK.

## FixedSizeRoundRobinReservoir
 
`FixedSizeRoundRobinReservoir` is an experimental reservoir that samples at most a fixed number of exemplars using a round-robin strategy to distribute measurements across independent buckets, each using Algorithm L for sampling.

This reservoir provides higher concurrent performance by sharding state across independent buckets. Dispatch rotates continuously across collection cycles, distributing measurements evenly across buckets. A sampling bias may be introduced if measurement patterns are periodic and coincide with the reservoir size.
