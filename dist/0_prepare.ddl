ALTER TABLE ExampleAlter ALTER COLUMN Col_1 INT64;
CREATE TABLE Tmp_ExampleRecreated (IdA INT64 NOT NULL, IdB INT64 NOT NULL, Meta STRING(MAX)) PRIMARY KEY (IdA, IdB), INTERLEAVE IN PARENT OutOfScope ON DELETE NO ACTION;
ALTER TABLE ExampleAlter ADD COLUMN Tmp_Col_3 INT64;