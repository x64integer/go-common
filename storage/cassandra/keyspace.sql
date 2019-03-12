CREATE KEYSPACE default_keyspace 
    WITH replication = {'class': 'NetworkTopologyStrategy', 'replication_factor' : 3, 'DC1': '3', 'DC2': '2'}