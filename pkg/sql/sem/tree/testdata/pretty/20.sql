CREATE TABLE t (a INT, b INT, c INT, PRIMARY KEY (a, b)) PARTITION BY LIST (a)  (
    PARTITION p1 VALUES IN (1),
    PARTITION p2 VALUES IN (2)
)
