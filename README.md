# dec
Fixed point decimal library

* Stored as int64 multiplied by 1billion (9 decimal places)
* All arithmetic operations except multiply and divide work as normal (use the Mul and Div functions)
* JSON/XML/SQL/gob serializable
* Convenience functions for common operations (avg, max, min, round, truncate etc)
