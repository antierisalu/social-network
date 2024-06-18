<script>
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
    console.log(fromUser)
</script>

<div class="message-container" {fromUser} {time} {msgID}>
    <div class="message-body">
        <a class="chat-username">{fromUsername}: </a>
        <a class="chat-message-content">{msgContent}</a>
        <div class="chat-time-wrapper">
            <a class="chat-time">{msgFormatedTime}</a>
        </div>
    </div>
</div>
