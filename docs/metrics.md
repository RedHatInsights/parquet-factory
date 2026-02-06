---
layout: page
title: Metrics
nav_order: 6
---


## Metrics

All the metrics are sent to our Prometheus Pushgateway] The name of the metrics and where they are gathered can be seen in the next schema:

<div class="mermaid">
flowchart LR;
    subgraph Metrics;
        m_offsets_marked(Offsets marked);
        m_offsets_consumed(Offsets consumed);
        m_offsets_processed(Offsets processed);
        m_files_generated(New files generated);
        m_inserted_rows(Inserted rows);
    end;
    subgraph Parquet Factory;
        PF1[Init] --> PF2[Connect to Kafka] --> PF3[Consume];

        PF3 --> m_offsets_marked;
        PF3 --> m_offsets_consumed;
        PF3 --> m_offsets_processed;

        PF3 --> PF5[Generate tables];

        PF5 --> m_files_generated;
        PF5 --> m_inserted_rows;
    end;
</div>

Apart from the ones listed in the schema, each one of these nodes will represent a state in the `state` metric. Being:


| State                  | Code |
|------------------------|------|
| Idle                   | 0    |
| Init                   | 1    |
| Connect to Kafka       | 2    |
| Consume                | 3    |
| Generate tables        | 4    |

So, the metrics are:

- `offsets_marked`: number of messages which offset has been [marked](https://pkg.go.dev/github.com/Shopify/sarama@v1.27.1?utm_source=gopls#PartitionOffsetManager.MarkOffset).
- `offsets_consummed`: number of messages [consumed](https://github.com/RedHatInsights/parquet-factory/-/blob/master/reportreader/reportreader.go#:~:text=c.limits.-,MessageProcessed,-()).
- `offsets_processed`: number of messages [processed](https://github.com/RedHatInsights/parquet-factory/-/blob/master/reportreader/reportreader.go#:~:text=c.offsetTracker.-,RecordOffset,-(m)).
- `files_generated`: number of files generated. Increased every time a file is [created](https://github.com/RedHatInsights/parquet-factory/-/blob/master/aggregator/rule_hit.go#:~:text=tracker.S3Writer.NewFile).
- `inserted_rows`: number of rows written ([check](https://github.com/RedHatInsights/parquet-factory/-/blob/master/parquet-factory.go#:~:text=tracker.WriteParquetFiles())).
- `state`: state of the cronjob.

There will be also an `error_count` metric.

<script src="{{ "/resources/mermaid.min.js" | relative_url }}"></script>
