<!-- REDUNDANT -->
<script>
    import Chatbox from "../chat/chatbox.svelte";
</script>
<footer>
    <Chatbox />
    <!-- <div class="chatBox">44444444444</div>
    <div class="chatBox">333333333333</div>
    <div class="chatBox">211111111111</div>
    <div class="chatBox">122222222222</div>
    <div class="searchDiv">
        <input type="search" placeholder="Search Chats">
    </div> -->
    
</footer>

<style>
    footer{
        /* border: 1px solid white; */
        height:100%;
        display:flex;
        flex-direction: row;
        justify-content: right;
    }
    .chatBox{
        color: #3c3c3c;
        font-size: 12px;
        border: rgb(141, 0, 141);

    }
    .searchDiv {
        width:220px;
        display: flex;
        justify-content: center; /* Horizontally center the input */
        align-items: center; /* Vertically center the input */
        height:100%;
        
    }
    input{
        margin: 0;
        border-color: #3c3c3c;
    }

</style>