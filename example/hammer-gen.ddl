ALTER TABLE ExampleAlter ALTER COLUMN Col_1 INT64;
ALTER TABLE ExampleAlter DROP COLUMN Col_3;
ALTER TABLE ExampleAlter ADD COLUMN Col_3 INT64;
ALTER TABLE ExampleAlter DROP COLUMN Col_2;
DROP TABLE ExampleRecreated;
CREATE TABLE ExampleRecreated (
  IdA INT64 NOT NULL,
  IdB INT64 NOT NULL,
  Meta STRING(MAX),
) PRIMARY KEY(IdA, IdB),
  INTERLEAVE IN PARENT OutOfScope ON DELETE NO ACTION;
