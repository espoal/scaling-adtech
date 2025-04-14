# scaling-adtech

In this repository there is an example of a simple ad server in golang. Due to time constraints, I had to make some
sacrifices, but I hope it will be helpful in giving a glimpse of my prioritization strategies. Given more time I would 
have:
- Implemented tests
- Spent more time on the validation
- Improved the documentation

Starts it with
```bash
docker compose up
```

Then you can check the health of the service with
```bash
curl http://localhost:8080/health
```

See the [OpenAPI document](api/openapi.yaml) for the complete API specification.


# Design decisions

The first design decision I took was to use Postgres as the database, as it fits well also with later design decisions.
Tightly coupled with this decision is the use of prepared statements: although not liked by everyone, they are a personal
favorite of mine, as they allow us to fully exploit the power of Postgres, and becomes the perfect key for our caching 
strategy later. Some people prefer to use an ORM, but in general I don't like indirections and find that ORMs turn 
databases like postgres in an ugly copy of mongodb, which I would pick instead if prepared statements are not an option.

The other important design decision is to use an event-driven architecture for the events tracking, a decision which will
again help us later in the scaling strategy. Using an antifragile message bus like Kafka or similar allow us not to 
worry about scalability issues, while setting the best foundation for our OLAP needs.

# AD rating

In the current design we rate ads simply by the bid price. The main hypothesis is that there will always be ads for 
a given placement and category, which is not trivial. With more time and with some historical data I would have 
used a machine learning approach: cluster categories and keywords according to relevance, increase the search space
to include adjacent categories and keywords, and then use a ranking algorithm to rate the ads. This would have 
allowed us to increase the average bid.

# Scaling considerations

We start with 2 instances on one node: each with 2 CPUs and 4GB of RAM, one for the API and one for the Postgres 
database. My guess is that the API will be the first bottleneck, so we will start with scaling it up to 48 CPUs.

I would expect the golang API service to scale almost linearly in this case, so soon enough postgres will become the 
bottleneck. We can then scale it up as well, but this will solve only part of the problem. The issue lies also in the 
database design: we query the table by btree (ad placement) AND by inclusion of a string in a vector
(category, keywords). This will cause performance issues as this type of queries is not well indexable by postgres,
as one query would require a btree index and the other a gin index. We could fix this by denormalizing the table and 
storing a reverse index by (placement, category), but this would make write more expensive, and would require careful
design considerations.

The next step, the classical one, would be to move postgres to a dedicated node, and use Redis for in memory caching
in the nodes running the API. Here the choice of using prepared statements comes in handy, as we can use them as the 
key for our Redis cache.

A more exotic approach which I prefer is to use RocksDB as a replica. While Redis uses ram, RocksDB uses SSDs, which 
greatly reduces costs while not sacrificing performance, as network will become the new bottleneck in this approach. 
Along with this we could reverse the caching problem, where each modification to our main postgres table will be 
published to a queue, and each instance of RocksDB will subscribe to it to keep the local cache updated. 

With the latter approach we could spawn many nodes to support the golang API and with a local RocksDB cache to scale out
operations. A single postgres cluster would be the source of truth for our operations, and I could eeasily envisions 
reaching millions of requests per second, rather than minute, with this approach.