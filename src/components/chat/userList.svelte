<script>
    import User from "../chat/user.svelte";

    // data
    let users = [
        {
            avatarPath: "./avatars/32.jpg",
            username: "Pepe TheFrog",
            userID: "32",
        },
        {
            avatarPath: "./avatars/32.jpg",
            username: "Pepe Lala",
            userID: "33",
        },
        {
            avatarPath: "./avatars/32.jpg",
            username: "Big Chungus",
            userID: "33",
        },
    ];

</script>

<div class="userListContainer">
    <div class="headerName">
        <h3>All Users</h3>
    </div>
    <div class="seperator"></div>

    <!-- Will contain all users to search from -->
    <div class="usersContainer">
        {#each users as user}
            <User avatarPath={user.avatarPath} username={user.username} userID={user.userID} />
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
        width: 100%;
        height: 100%;
    }

    .usersContainer {
        width: 100%;
        height: 95%;
    }

    .searchBarWrapper {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        width: 100%;
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