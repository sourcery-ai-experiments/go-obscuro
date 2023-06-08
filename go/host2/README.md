The host is the gatekeeper of the obscuro enclaves.

It provides an access bridge between the enclave and the outside world, in particular to the peer nodes in the Obscuro 
network, the L1 and user requests.

## Services
The host is organised as a set of services, each of which is responsible for a particular aspect of the host's responsibilities.

Services will share a common interface including Stop/Start, a health status and recent metrics. The host records a log of service statuses/stats for debugging.

Services can depend on other services, but they will request the service on each usage via a HostServices interface which can return an error if that service is in a bad state. 

## Enclaves
The most important services are the 'enclave-guardian' services which are responsible for monitoring the state of the enclaves and providing them with whatever they need.

A host may be responsible for multiple enclaves for HA purposes. In this case, the host will have multiple enclave-guardian services, one for each enclave, independently feeding blocks and providing two-way communication.

A sequencer host will have a primary enclave at any given time which is responsible for producing the batches and rollups. The other enclaves are hot back-ups.
