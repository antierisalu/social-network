<script>
    import MsgNotification from "../icons/msgNotification.svelte";
    export let avatarPath = "";
    export let username = "";
    export let userID = "";
</script>

<div class="user">
    <div class="profilePictureWrapper">
        <img src={avatarPath} alt={userID}>
    </div>

    <div class="usernameWrapper">
        <h2 class="username" style="margin: 0;">{username}</h2>
    </div>

    <div class="messageNotification">
        <!-- Slide/Appear svelte animation ***TODO -->
        <MsgNotification />
    </div>
</div>

<style>
.user {
        display: flex;
        justify-content: center;
        align-items: center;
        margin: 3%;
        width: 94%;
        height: 42px;
        border-radius: 5px;
        border: 1px solid rgb(145, 145, 145);
    }

    .profilePictureWrapper {
        margin-left: 6px;
        margin-top: 3px;
        margin-bottom: 3px;
        width: 34px;
        height: 34px;
        border-radius: 50%;
        border: 2px solid #699b1e;
    }
    
    .profilePictureWrapper img {
        width: 100%;
        height: 100%;
        object-fit: cover;
        border-radius: 50%;
    }

    .usernameWrapper {
        display: flex;
        justify-content: center;
        align-items: center;
        flex: 7;
        height: 100%;
    }

    .username {
        font-size: medium;
    }

    .messageNotification {
        flex: 1;
        margin-right: 6px;
        margin-top: 3px;
        margin-bottom: 3px;
        height: 100%;
        display: flex;
        justify-content: center;
        align-items: center;
    }
</style>