# DEMO of Type2

this is a demo of type2 AID

## Concept

a todo list can use aid to log in a remote todo list application, which can share the todo list with others.

![image](../../doc/overview-t2.png)

2 extra buttons: `bug` and `sign` can try the feature of the type2 aid demo.

## How to execute

1. start the AID server (use `make` in `(root)/aid-server`)
2. start the backend (use `make`)
3. start the frontend (use `yarn start`)

## Note

- you can clear DB in `(root)/aid-server` by `make clean`
- you also need to clear browser local storage to do experiments again
- after use `bug` button register aid in AID Server, although you can login with the same aid, you can't share the todo list with others, because your cert can not pass verification from aid server.
- backend can also forbid the login that can't pass the verification of the aid server.