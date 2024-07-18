# TERE POISID

Node must be installed

```bash
# After you install the repo run this (needed only once)
npm install
```

```bash
# After doing svelte updates you gotta compile it with
npm run dev
```

[Exercise](https://github.com/01-edu/public/tree/master/subjects/social-network)

[Audit Req](https://github.com/01-edu/public/tree/master/subjects/social-network/audit)

[Eraser](https://app.eraser.io/workspace/CJOPbqoi4KGm18qrX1qE)

[Trello](https://trello.com/b/JhSNJVWG/social-network)

image upload on register does not work

When you cancel follow request other person can still accept the request and you end up following that person anyway.

cannot close invite user to the group dropdown

Event created, which is less than 2 hours from the current time does not show up

Group members are not receiving notification about the new event created the group


##### Open two browsers, log in with different users on each browser, become part of the same group with both users and with one of the users create an event.

###### Did the other user received a notification regarding the creation of the event?

#### Docker

##### Try to run the application and use the docker command `"docker ps -a"`

###### Can you confirm that there are two containers (backend and frontend), and both containers have non-zero sizes indicating that they are not empty?

##### Try to access the social network application through your web browser.

###### Were you able to access the social network application through your web browser after running the docker containers, confirming that the containers are running and serving the application as expected?

#### Bonus

###### +Can you log in using Github or other type of external OAuthenticator (open standard for access delegation)?

###### +Did the student created a migration to fill the database?

###### +If you unfollow a user, do you get a confirmation pop-up?

###### +If you change your profile from public to private (or vice versa), do you get a confirmation pop-up?

###### +Is there other notification apart from the ones explicit on the subject?

###### +Does the project present a script to build the images and containers? (using a script to simplify the build)

###### +Do you think in general this project is well done?