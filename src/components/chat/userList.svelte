<script>
    import User from "../chat/user.svelte";
    import { allUsers, userInfo, onlineUserStore, lastMsgStore, chatNotifStore } from "../../stores";
    $: users = $allUsers;
    $: onlineUsers = $onlineUserStore
    $: lastMsgs = $lastMsgStore
    // $: notification = $chatNotifStore
    // Reactive declaration for filtered users (searchBar)
    var searchQuery = "";
    $: filteredUsers = searchQuery ? searchUsers(searchQuery) : sortedUsers;


    $: sortedUsers = [...users].sort((a, b) => {
        // 1. Sort by last message timestamps
        const aLastMsg = lastMsgs[a.ID] ? new Date(lastMsgs[a.ID]) : new Date(0);
        const bLastMsg = lastMsgs[b.ID] ? new Date(lastMsgs[b.ID]) : new Date(0);
        if (aLastMsg > bLastMsg) {
            return -1;
        } else if (aLastMsg < bLastMsg) {
            return 1;
        }
    
        // 2. Sort by online status
        const aISonline = onlineUsers.includes(a.ID);
        const bISonline = onlineUsers.includes(b.ID);
        if (aISonline && !bISonline) {
            return -1;
        } else if (!aISonline && bISonline) {
            return 1;
        } else {
            // 3. incase same status sort by first name
            return a.FirstName.localeCompare(b.FirstName);
        }
    });

    // Stole this from searchBar.svelte
    const searchUsers = (searchQuery) => {
        const [firstNameQuery, lastNameQuery] = searchQuery
        .toLowerCase()
        .trim()
        .split(" ");
        return users.filter((user) => {
            if (!firstNameQuery && !lastNameQuery) {
                return;
            }

            let firstNameMatch;
            let firstNameMatch2;
            let lastNameMatch;
            let lastNameMatch2;

            if (!lastNameQuery || (firstNameQuery && lastNameQuery)) {
                firstNameMatch = user.FirstName.toLowerCase().includes(firstNameQuery);
                firstNameMatch2 = user.LastName.toLowerCase().includes(firstNameQuery);
            }

            if (!firstNameQuery || (firstNameQuery && lastNameQuery)) {
                lastNameMatch = user.LastName.toLowerCase().includes(lastNameQuery);
                lastNameMatch2 = user.FirstName.toLowerCase().includes(lastNameQuery);
            }

            return (
                firstNameMatch || lastNameMatch || firstNameMatch2 || lastNameMatch2
            );
        });
    };
</script>

<div class="userListContainer">
    <div class="headerName">
        <h3>All Users</h3>
    </div>
    <div class="seperator"></div>

    <!-- Will contain all users to search from -->
    <div class="usersContainer" id="usersContainer">
        <!-- {#each sortedUsers as user} -->
        {#each filteredUsers as user}
            {#if user.ID !== $userInfo.id}
                <User 
                avatarPath={user.Avatar} 
                firstName={user.FirstName} 
                lastName={user.LastName}
                userID={user.ID}
                isOnline={onlineUsers.includes(user.ID)}
                />
            {/if}
        {/each}
    </div>

    <div class="seperator"></div>
    <div class="searchBarWrapper">
        <input class="searchBar" type="search" bind:value={searchQuery} placeholder="Search Chats">
    </div>
</div>

<style>
    h3 {
        color: white;
        margin: 0
    }

    .seperator {
        margin-top: 5px;
        margin-bottom: 5px;
        background-color: rgb(119, 119, 119);
        border-radius: 7px;
        height: 3px;
    }
    
    .userListContainer {
        display: flex;
        flex-direction: column;
        justify-content: center;
        width: 100%;
        height: 100%;
        max-height: 100%;
    }

    .usersContainer {
        width: 100%;
        max-height: 95%;
        height: 95%;
        overflow-y: scroll;
        scrollbar-width: thin;
        scrollbar-color:  greenyellow #011;
    }

    .searchBarWrapper {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        width: 100%;
        /* max-height:  */
        height: 5%;
        position:relative;
        bottom: 0px;
    }

    .searchBar {
        margin: 0;
        width: 100%;
        height: 100%;
        border-top: none;
    }
</style>