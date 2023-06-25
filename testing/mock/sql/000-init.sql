USE staff_api_test;

DROP TABLE IF EXISTS luckperms_players;
DROP TABLE IF EXISTS luckperms_user_permissions;

CREATE TABLE luckperms_players
(
    uuid          varchar(36) NOT NULL,
    username      varchar(16) NOT NULL,
    primary_group varchar(16) NOT NULL,
    PRIMARY KEY (uuid)
);

CREATE TABLE luckperms_user_permissions
(
    id int(11) NOT NULL AUTO_INCREMENT,
    uuid varchar(36) NOT NULL,
    permission varchar(200),
    value tinyint(1),
    server varchar(36),
    world varchar(36),
    expiry int(11),
    contexts varchar(200),
    PRIMARY KEY (id)
);

CREATE INDEX luckperms_players_uuid ON luckperms_players (uuid);
CREATE INDEX luckperms_groups_uuid ON luckperms_user_permissions (uuid);

# Insert players with random UUIDs
INSERT INTO luckperms_players (uuid, username, primary_group)
VALUES ('40a1b924-e5a8-4444-9fec-73db15ee7c8d', 'player1', 'test1');
INSERT INTO luckperms_players (uuid, username, primary_group)
VALUES ('c806a034-b626-4336-a1e2-37b902888bf5', 'player2', 'test1');
INSERT INTO luckperms_players (uuid, username, primary_group)
VALUES ('d0b765a7-87d6-49fb-96d4-17e1cdb6ca2e', 'player3', 'test1');
INSERT INTO luckperms_players (uuid, username, primary_group)
VALUES ('c706ec9f-1d29-4cf6-a849-ba5830b5cb41', 'player4', 'test2');
INSERT INTO luckperms_players (uuid, username, primary_group)
VALUES ('b1a2291f-3955-431e-b274-2388f85d3b63', 'player5', 'test2');
INSERT INTO luckperms_players (uuid, username, primary_group)
VALUES ('ef7eb665-a3ac-40a6-b9a4-1100f60b28cd', 'player6', 'test3');

# Insert permissions
INSERT INTO luckperms_user_permissions (uuid, permission, value)
VALUES ('40a1b924-e5a8-4444-9fec-73db15ee7c8d', 'group.test1', 1);
INSERT INTO luckperms_user_permissions (uuid, permission, value)
VALUES ('c806a034-b626-4336-a1e2-37b902888bf5', 'group.test1', 1);
INSERT INTO luckperms_user_permissions (uuid, permission, value)
VALUES ('d0b765a7-87d6-49fb-96d4-17e1cdb6ca2e', 'group.test1', 1);
INSERT INTO luckperms_user_permissions (uuid, permission, value)
VALUES ('c706ec9f-1d29-4cf6-a849-ba5830b5cb41', 'group.test2', 1);
INSERT INTO luckperms_user_permissions (uuid, permission, value)
VALUES ('b1a2291f-3955-431e-b274-2388f85d3b63', 'group.test2', 1);
INSERT INTO luckperms_user_permissions (uuid, permission, value)
VALUES ('ef7eb665-a3ac-40a6-b9a4-1100f60b28cd', 'group.test3', 1);

