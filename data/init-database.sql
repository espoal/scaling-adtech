CREATE TYPE line_item_status AS enum ('active', 'paused', 'completed');

CREATE TABLE line_items (
    id           varchar(255) PRIMARY KEY,
	name         varchar(255),
	advertiser_id varchar(255),
	bid          double precision,
	budget       double precision,
	placement    varchar(255),
	categories   varchar(255)[],
	keywords     varchar(255)[],
	status       line_item_status,
	created_at    timestamp,
	updated_at    timestamp
);

CREATE INDEX placement_index ON line_items (placement);
CREATE INDEX placement_index_with_category ON line_items (placement) include (categories);
CREATE INDEX placement_index_with_category_and_keywords ON line_items (placement) include (categories, keywords);
