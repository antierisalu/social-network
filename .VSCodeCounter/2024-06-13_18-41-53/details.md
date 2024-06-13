# Details

Date : 2024-06-13 18:41:53

Directory /home/nitroarch/kood/social-network

Total : 51 files,  3476 codes, 116 comments, 471 blanks, all 4063 lines

[Summary](results.md) / Details / [Diff Summary](diff.md) / [Diff Details](diff-details.md)

## Files
| filename | language | code | comment | blank | total |
| :--- | :--- | ---: | ---: | ---: | ---: |
| [README.md](/README.md) | Markdown | 14 | 0 | 8 | 22 |
| [backend/go.mod](/backend/go.mod) | Go Module File | 15 | 0 | 5 | 20 |
| [backend/go.sum](/backend/go.sum) | Go Checksum File | 29 | 0 | 1 | 30 |
| [backend/pkg/auth.go](/backend/pkg/auth.go) | Go | 231 | 36 | 36 | 303 |
| [backend/pkg/avatars.go](/backend/pkg/avatars.go) | Go | 99 | 0 | 18 | 117 |
| [backend/pkg/db/migrations/sqlite/1_create_userstable.down.sql](/backend/pkg/db/migrations/sqlite/1_create_userstable.down.sql) | SQL | 1 | 0 | 0 | 1 |
| [backend/pkg/db/migrations/sqlite/1_create_userstable.up.sql](/backend/pkg/db/migrations/sqlite/1_create_userstable.up.sql) | SQL | 16 | 0 | 4 | 20 |
| [backend/pkg/db/migrations/sqlite/2_create_followerstable.down.sql](/backend/pkg/db/migrations/sqlite/2_create_followerstable.down.sql) | SQL | 1 | 0 | 0 | 1 |
| [backend/pkg/db/migrations/sqlite/2_create_followerstable.up.sql](/backend/pkg/db/migrations/sqlite/2_create_followerstable.up.sql) | SQL | 9 | 0 | 3 | 12 |
| [backend/pkg/db/migrations/sqlite/3_create_groups_and_groupmembers.down.sql](/backend/pkg/db/migrations/sqlite/3_create_groups_and_groupmembers.down.sql) | SQL | 2 | 0 | 1 | 3 |
| [backend/pkg/db/migrations/sqlite/3_create_groups_and_groupmembers.up.sql](/backend/pkg/db/migrations/sqlite/3_create_groups_and_groupmembers.up.sql) | SQL | 16 | 0 | 4 | 20 |
| [backend/pkg/db/migrations/sqlite/4_create_rest_of_db.down.sql](/backend/pkg/db/migrations/sqlite/4_create_rest_of_db.down.sql) | SQL | 8 | 0 | 0 | 8 |
| [backend/pkg/db/migrations/sqlite/4_create_rest_of_db.up.sql](/backend/pkg/db/migrations/sqlite/4_create_rest_of_db.up.sql) | SQL | 77 | 0 | 7 | 84 |
| [backend/pkg/db/sqlite/database.db](/backend/pkg/db/sqlite/database.db) | SQL | 63 | 0 | 21 | 84 |
| [backend/pkg/db/sqlite/sqlite.go](/backend/pkg/db/sqlite/sqlite.go) | Go | 27 | 1 | 8 | 36 |
| [backend/pkg/profile.go](/backend/pkg/profile.go) | Go | 141 | 8 | 28 | 177 |
| [backend/pkg/search.go](/backend/pkg/search.go) | Go | 51 | 4 | 9 | 64 |
| [backend/pkg/structs.go](/backend/pkg/structs.go) | Go | 33 | 2 | 6 | 41 |
| [backend/pkg/websocket.go](/backend/pkg/websocket.go) | Go | 69 | 1 | 15 | 85 |
| [backend/server.go](/backend/server.go) | Go | 23 | 4 | 8 | 35 |
| [go.mod](/go.mod) | Go Module File | 2 | 0 | 2 | 4 |
| [go.sum](/go.sum) | Go Checksum File | 0 | 0 | 1 | 1 |
| [package-lock.json](/package-lock.json) | JSON | 972 | 0 | 1 | 973 |
| [package.json](/package.json) | JSON | 24 | 0 | 1 | 25 |
| [public/global.css](/public/global.css) | CSS | 54 | 0 | 12 | 66 |
| [public/index.html](/public/index.html) | HTML | 14 | 0 | 5 | 19 |
| [rollup.config.js](/rollup.config.js) | JavaScript | 54 | 14 | 11 | 79 |
| [scripts/setupTypeScript.js](/scripts/setupTypeScript.js) | JavaScript | 78 | 31 | 26 | 135 |
| [src/App.svelte](/src/App.svelte) | Svelte | 52 | 0 | 7 | 59 |
| [src/components/auth/login.svelte](/src/components/auth/login.svelte) | Svelte | 102 | 0 | 10 | 112 |
| [src/components/auth/register.svelte](/src/components/auth/register.svelte) | Svelte | 188 | 0 | 23 | 211 |
| [src/components/chat/chat.svelte](/src/components/chat/chat.svelte) | Svelte | 23 | 0 | 7 | 30 |
| [src/components/chat/user.svelte](/src/components/chat/user.svelte) | Svelte | 70 | 1 | 9 | 80 |
| [src/components/chat/userList.svelte](/src/components/chat/userList.svelte) | Svelte | 63 | 1 | 9 | 73 |
| [src/components/groups/groups.svelte](/src/components/groups/groups.svelte) | Svelte | 14 | 0 | 2 | 16 |
| [src/components/icons/msgNotification.svelte](/src/components/icons/msgNotification.svelte) | Svelte | 13 | 0 | 7 | 20 |
| [src/components/notifications/notifications.svelte](/src/components/notifications/notifications.svelte) | Svelte | 10 | 0 | 2 | 12 |
| [src/components/profile/privateData.svelte](/src/components/profile/privateData.svelte) | Svelte | 104 | 0 | 17 | 121 |
| [src/components/profile/profile.svelte](/src/components/profile/profile.svelte) | Svelte | 147 | 0 | 28 | 175 |
| [src/components/profile/searchBar.svelte](/src/components/profile/searchBar.svelte) | Svelte | 86 | 2 | 15 | 103 |
| [src/components/structure/footer.svelte](/src/components/structure/footer.svelte) | Svelte | 33 | 0 | 3 | 36 |
| [src/components/structure/header.svelte](/src/components/structure/header.svelte) | Svelte | 67 | 0 | 7 | 74 |
| [src/components/structure/mainpage.svelte](/src/components/structure/mainpage.svelte) | Svelte | 60 | 4 | 12 | 76 |
| [src/components/structure/mainwindow.svelte](/src/components/structure/mainwindow.svelte) | Svelte | 75 | 1 | 12 | 88 |
| [src/main.js](/src/main.js) | JavaScript | 7 | 0 | 2 | 9 |
| [src/shared/button.svelte](/src/shared/button.svelte) | Svelte | 52 | 0 | 9 | 61 |
| [src/shared/imagePreview.svelte](/src/shared/imagePreview.svelte) | Svelte | 98 | 0 | 14 | 112 |
| [src/shared/matrix.svelte](/src/shared/matrix.svelte) | Svelte | 43 | 0 | 15 | 58 |
| [src/stores.js](/src/stores.js) | JavaScript | 15 | 6 | 10 | 31 |
| [src/utils.js](/src/utils.js) | JavaScript | 6 | 0 | 2 | 8 |
| [src/websocket.js](/src/websocket.js) | JavaScript | 25 | 0 | 8 | 33 |

[Summary](results.md) / Details / [Diff Summary](diff.md) / [Diff Details](diff-details.md)