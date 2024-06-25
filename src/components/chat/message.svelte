<script>
    import {userInfo} from "../../stores";
    $: user = $userInfo.id;
    export let fromUser;
    export let fromUsername;
    export let time;
    export let msgID;
    export let msgContent;

    // Formats 2024-04-30 20:11:27 ---> 30/04/24' @20:11
    function formatChatDateTime(timeString) {
        const dateTime = new Date(timeString);
        let day = dateTime.getDate();
        let month = dateTime.getMonth() + 1;
        let year = dateTime.getFullYear() % 100;
        let hours = dateTime.getHours();
        let minutes = dateTime.getMinutes();
        day = (day < 10 ? '0' : '') + day;
        month = (month < 10 ? '0' : '') + month;
        year = (year < 10 ? '0' : '') + year;
        hours = (hours < 10 ? '0' : '') + hours; 
        minutes = (minutes < 10 ? '0' : '') + minutes;
        const formatted = `${day}/${month}/${year}' @${hours}:${minutes}`
        return formatted
    }
    let msgFormatedTime = formatChatDateTime(time);

// teen checki ehk kui user == fromUser 
// kui on sama siis anna talle class > style > peida kasutaja nimi ja liiguta tekts paremale
//else anna class > style > jne

$: if (user) {
    console.log('testkaka:', user, fromUser);
  }

</script>

<div class="message-container" {fromUser} {time} {msgID}>
    <div class="message-body">
        {#if user == fromUser}
        <div class="chat-username-owner">{fromUsername}: </div>
        {:else}
        <div class="chat-username-quest">{fromUsername}: </div>
        {/if}
        {#if user == fromUser}
        <div class="chat-message-content-owner">{msgContent}</div>
        {:else}
        <div class="chat-message-content-quest">{msgContent}</div>
        {/if}
        <div class="chat-time-wrapper">
            <div class="chat-time">{msgFormatedTime}</div>
        </div>
    </div>
</div>

<style>

:global(.message-container) {
        display: flex;
        flex-direction: column;
    }

    :global(.message-body) {
        display: flex;
        flex-direction: column;
        width: 100%;
        margin-bottom: 15px;
    }
   :global(.chat-time-wrapper) {
        width: 100%;
        display: flex;
        justify-content: end;
        margin-top: 2px;
    }

    :global(.chat-username-owner) {
        display: none;
    }

    :global(.chat-username-quest) {
        font-size: large;
        color: white;
        font-weight: 700;
        min-height: 18px;
        user-select: none;
        text-align: left;
        margin-left: 10px;
        margin-bottom: 50px;
    }
    :global(.chat-message-content-owner) {
        font-size: medium;
        color: black;
        font-weight: 600;
        text-align: right;
        margin-right: 20px;
        background-color: blue;
        border-radius: 15px;
        padding: 10px;
        word-wrap: break-word;
        align-self: flex-end;
        
    }

    :global(.chat-message-content-quest) {
        font-size: medium;
        color: black;
        font-weight: 600;
        text-align: left;
        margin-left: 20px;
        background-color: lightblue;
        border-radius: 15px;
        padding: 10px;
        word-wrap: break-word;
        align-self: flex-start;
    }

    :global(.message-container:hover .chat-time) {
        top: 0;
        opacity: 1;
    }

    :global(.chat-time) {
        position: relative;
        color: gray;
        font-weight: 500;
        min-height: 18px;
        user-select: none;
        opacity: 0;
        transition: top 0.3s ease, opacity 0.3s ease;
    }
    
</style>
