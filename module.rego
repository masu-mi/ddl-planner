package ddl
import future.keywords.in

alter_table[x] {
	input[x].Alter
}
drop_table[x] {
	input[x].Drop
}
create_table[x] {
	input[x].Create
}

alter_drop_col[id] {
	input[id].TableAlteration.Drop
	alter_table[id]
}

alter_add_col[id] {
	input[id].TableAlteration.Add
	alter_table[id]
}

add[c_id] { data.ddl.create_table[c_id] }
add[a_id] { some a_id in { id | data.ddl.alter_table[id] } - { id | data.ddl.alter_drop_col[id] } }

del[id] { data.ddl.drop_table[id] }
del[id] { data.ddl.alter_drop_col[id] }
