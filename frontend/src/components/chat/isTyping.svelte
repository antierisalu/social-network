<script>
    export let isTyping;
    export let userName;
    export let isGroup;
    export let chatID;
    let self;
    let multipleTyping = false;
    const maxBongos = 6;
    import { groupIsTypingStore } from "../../stores";
    $: groupTypingStore = $groupIsTypingStore
    let audioMap = new Map();
    $: {
        let containerElem = document.querySelector(`.chatBox[userid="${"GroupChatID_"+chatID}"]`);
        if (groupTypingStore[chatID] !== undefined){
            const currentTyping = groupTypingStore[chatID];
            // Start audio for new typers
            currentTyping.forEach(element => {
                if (!audioMap.has(element)) {
                    let audio = new Audio("typing.mp3");
                    audio.volume = 0.01; // 1% volume, DO NOT INCREASE
                    audio.loop = true;
                    audio.play();
                    audioMap.set(element, audio);
                }
            });
            // Stop audio for people who stopped typing
            audioMap.forEach((audio, person) => {
                if (!currentTyping.includes(person) || !containerElem) {
                    audio.pause();
                    audioMap.delete(person);
                }
            });
        }
    }
    $: {
        if (isGroup) {
            if (groupTypingStore[chatID] === undefined) {
                isTyping = false
            } else if (groupTypingStore[chatID].length === 0) {
                isTyping = false;
            } else if (groupTypingStore[chatID].length === 1){
                isTyping = true;
                userName = groupTypingStore[chatID][0]
            } else {
                isTyping = true
                multipleTyping = true;
                userName = "multiple people"
            }
        }
    }
    $: bongoIndexes = [...Array(Math.min(groupTypingStore[chatID]?.length || 0, maxBongos)).keys()];
</script>

<div bind:this={self} id="isTypingWrapper" style="{isTyping ? 'display:block' : 'display:none'}">
    <div class="glowContainer typingGlow"></div>
    {#if multipleTyping}
        {#each bongoIndexes as index}
            <img id="catO"src="bongo-cat-transparent.webp" alt="" style="translate: {index * 3}px; z-index: {6 + index};">
        {/each}
    {:else}
        <img id="catO"src="bongo-cat-transparent.webp" alt="">
    {/if}
    <div class="container">
        <p>{userName} {#if multipleTyping}are{:else}is{/if} typing</p>
        <div class="typingAnimation">
            <div class="circle c01"></div>
            <div class="circle c02"></div>
            <div class="circle c03"></div>
        </div>
    </div>
</div>

<style>
    .container {
        margin-left: 68px;
        height: 91%;
        display: flex;
        align-items: center;
        justify-content: left;
    }

    .container p {
        user-select: none;
        font-size: small;
        font-weight: 500;
        text-align: bottom;
        margin: 0;
        padding: 0;
        bottom: 32px;
        width: max-content;
        height: max-content;
    }

    .glowContainer {
        position: absolute;
        bottom: 2px;
        width: 100%;
        height: 100%;
        border-bottom-left-radius: 10px;
        border-bottom-right-radius: 10px;   
    }

    #catO {
        z-index: 4;
        position: absolute;
        width: 50px;
        -webkit-transfrom: scaleX(-1);
        transform: scaleX(-1);
        rotate: 14deg;
        bottom: 32px;
        left: 16px;
    }

    #isTypingWrapper {
        position: absolute;
        bottom: 2px;
        width: 100%;
        height: 100px;
        border-bottom-left-radius: 10px;
        border-bottom-right-radius: 10px;

    }

    :global(.typingGlow) {
        z-index: -1;
        display: block;
        background: linear-gradient(0deg, rgba(138, 200, 221, 0.714) 43%, rgba(178,4,254,0) 92%); 
        animation: pulseGlow 1.5s infinite;
    }
    
    @keyframes pulseGlow {
        0% {
            opacity: 0.2;
        }
        50% {
            opacity: 1;
        }
        100% {
            opacity: 0;
        }
    }

    .typingAnimation {
        transform: translate(-1px, -1px);
        display: flex;
        flex-direction: row;
        justify-content: space-evenly;
        align-items: end;
        width: 20px;
        height: 14px;
    }

    .circle {
        background: rgb(153, 153, 153);
        width: 3px;
        height: 3px;
        border-radius: 100%;
        animation: wave 1.5s infinite;
    }
    .c01 {
        animation-delay: 0.1s;
    }
    .c02 {
        animation-delay: 0.2s;
    }
    .c03 {
        animation-delay: 0.3s;
    }
    
    @keyframes wave {
        0% {
            transform: translateY(-0.5px);
        }
        45% {
            transform: translateY(-4px);
        }
        100% {
            transform: translateY(0px);
        }
    }
</style>