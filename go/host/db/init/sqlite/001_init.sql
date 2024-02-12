create table if not exists rollup
(
    hash              binary(16) primary key,
    start_seq         int        NOT NULL,
    end_seq           int        NOT NULL,
    started_at        int        NOT NULL,
    compression_block binary(16) NOT NULL
);

create index IDX_ROLLUP_PROOF on rollup (compression_block);
create index IDX_ROLLUP_SEQ on rollup (start_seq, end_seq);

create table if not exists batch_body
(
    id          int        NOT NULL primary key,
    content     mediumblob NOT NULL
);

create table if not exists batch
(
    sequence       int primary key,
    hash           binary(16) NOT NULL unique,
    height         int        NOT NULL,
    tx_count       int        NOT NULL,
    header         blob       NOT NULL,
    body           int        NOT NULL REFERENCES batch_body
    );
create index IDX_BATCH_HASH on batch (hash);
create index IDX_BATCH_HEIGHT on batch (height);