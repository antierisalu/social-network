<script>
    import { userInfo, IMAGE_URL } from "../../stores";
    import { fade } from "svelte/transition";
    $: user = $userInfo.id;
    export let fromUser;
    export let fromUsername;
    export let time;
    export let msgID;
    export let msgContent;
    export let AvatarPath = "";
    if (AvatarPath === "") {
        AvatarPath = "/avatars/default.png";
    }

    // Formats 2024-04-30 20:11:27 ---> 30/04/24' @20:11
    function formatChatDateTime(timeString) {
        const dateTime = new Date(timeString);
        let day = dateTime.getDate();
        let month = dateTime.getMonth() + 1;
        let year = dateTime.getFullYear() % 100;
        let hours = dateTime.getHours();
        let minutes = dateTime.getMinutes();
        day = (day < 10 ? "0" : "") + day;
        month = (month < 10 ? "0" : "") + month;
        year = (year < 10 ? "0" : "") + year;
        hours = (hours < 10 ? "0" : "") + hours;
        minutes = (minutes < 10 ? "0" : "") + minutes;
        const formatted = `${day}/${month}/${year}' @${hours}:${minutes}`;
        return formatted;
    }
    let msgFormatedTime = formatChatDateTime(time);
    $: showTime = false;
</script>

<div class="message-container" {fromUser} {time} {msgID}>
    {#if user == fromUser}
        <!-- svelte-ignore a11y-mouse-events-have-key-events -->
        <div
            class="chat-message-content-owner"
            on:mouseover={() => {
                showTime = true;
            }}
            on:mouseout={() => {
                showTime = false;
            }}
        >
            {msgContent}
        </div>
        {#if showTime === true}
            <div
                class="chat-time"
                transition:fade={{ delay: 1000, duration: 500 }}
            >
                {msgFormatedTime}
            </div>
        {/if}
    {:else}
        <div class="chat-message-quest">
            <img src="{IMAGE_URL}{AvatarPath}" alt="userID" />
            <!-- svelte-ignore a11y-mouse-events-have-key-events -->
            <div
                class="chat-message-content-quest"
                on:mouseover={() => {
                    showTime = true;
                }}
                on:mouseout={() => {
                    showTime = false;
                }}
            >
                {msgContent}
            </div>
            {#if showTime === true}
                <div
                    class="chat-time"
                    transition:fade={{ delay: 1000, duration: 500 }}
                >
                    {fromUsername}
                    {msgFormatedTime}
                </div>
            {/if}
        </div>
    {/if}
</div>

<style>
    .message-container {
        display: flex;
        flex-direction: column;
        width: 100%;
        margin-bottom: 15px;
        position: relative;
    }

    .chat-message-content-owner {
        max-width: 180px;
        font-size: medium;
        color: white;
        font-weight: 600;
        text-align: right;
        margin-right: 20px;
        background-color: blue;
        border-radius: 15px;
        padding: 10px;
        word-wrap: break-word;
        align-self: flex-end;
    }
    .chat-message-quest {
        display: flex;
        align-items: center;
    }
    .chat-message-content-quest {
        max-width: 180px;
        font-size: medium;
        color: black;
        font-weight: 600;
        text-align: left;
        margin-left: 5px;
        background-color: lightblue;
        border-radius: 15px;
        padding: 10px;
        word-wrap: break-word;
        align-self: flex-start;
    }

    img {
        width: 20px;
        height: 20px;
        object-fit: cover;
        border-radius: 50%;
    }
    .chat-time {
        position: absolute;
        top: 0;
        white-space: nowrap;
        background-color: #f0f0f0;
        padding: 2px 5px;
        border-radius: 10px;
        font-size: 0.75rem;
        color: #333;
    }
</style>
