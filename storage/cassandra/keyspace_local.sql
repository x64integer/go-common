CREATE KEYSPACE default_keyspace 
    WITH replication = {'class': 'SimpleStrategy', 'replication_factor' : 1}