<script>
    import User from "../chat/user.svelte";
    import { allUsers } from "../../stores";
    $: users = $allUsers;
</script>

<div class="userListContainer">
    <div class="headerName">
        <h3>All Users</h3>
    </div>
    <div class="seperator"></div>

    <!-- Will contain all users to search from -->
    <div class="usersContainer">
        <!-- TODO: IF user is current user dont display them in the list (Incase we dont want user controlled
        status states 'online'/'busy' toggle and want only 'online/offline') ***-->
        {#each users as user}
            <User 
            avatarPath={user.Avatar} 
            firstName={user.FirstName} 
            lastName={user.LastName}
            userID={user.ID} />
        {/each}
    </div>

    <div class="seperator"></div>
    <div class="searchBarWrapper">
        <input class="searchBar" type="search" placeholder="Search Chats">
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