package process

/*
`process` should be the handler of the server layer, receving collection result from SDK and send respective payload to
storage. (MongoDB, ElasticSearch, and MinIO)
`process` layer should also propogate configuration update to each client SDK. The structure is in a node/cluster there
is one agent connecting to N client SDK, and the cloud connects to M node agent, a hierachical structure. In this way,
load on the upper layer can be easier to handle, compared to cloud:SDK = 1:NxM
*/
