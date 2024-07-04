<script>
    import { fade } from "svelte/transition";
    import Button2 from "../../shared/button2.svelte";
    import { userInfo, userProfileData, allUsers} from "../../stores";
    import { sendFollow, selectUser} from "../../utils";
    import { sendMessage } from "../../websocket";
    export let userID;
    export let userName;
    export let user;

    function handleClick(action, target) {
        // Limiter for request spam *DISABLED ATM* (timeout val = 0)
        setTimeout(() => {
            sendFollow(action, target.ID)
            if (action === 1) {
                var messageData = {
                    type: "followNotif",
                    targetid: target.ID,
                    fromid: $userInfo.id,
                    data: String
                }
                messageData.data = "follow_" + (target.ID).toString()
                sendMessage(JSON.stringify(messageData))
            } else if (action === 0) {
                var messageData = {
                    type: "followRequestNotif",
                    targetid: target.ID,
                    fromid: $userInfo.id,
                    data: String
                }
                messageData.data = "followRequest_" + (target.ID).toString()
                sendMessage(JSON.stringify(messageData))
            } else if (action === -1) {
                // handle cancel request ***TODO
            }
            if ($userProfileData.id === target.ID) {
                $userProfileData.areFollowing = action
                $userProfileData = $userProfileData
            }
        },0); // was 300, but likely can refactor setTimeout out later. ***TODO
    }
</script>

<!-- THIS DOM IS HALF-BAKED ***TODO: Continue Chris -->
<div class="container" {userID} {userName}>
    <div class="sub-container">
        <!-- <p>Follow {userName} to start a chat</p> -->
        <div class="buttons">
            {#if user.AreFollowing == 0}
                <div in:fade>
                    <Button2 
                        btnText="Cancel Request"
                        onClick={() => handleClick(-1, user)}
                        styleConfig={{ 
                            btnType: 2, 
                            btnWidth: 140,
                            btnHeight: 42,
                            fontSize: 15, 
                            fontColor: "#011",
                            fontWeight: "450",
                            hoverFontColor: "crimson",
                            borderRadius: 10,
                            borderWidth: 2,
                            borderColor: "crimson",
                            borderHoverColor: "red",
                            backgroundTone1: "crimson", 
                            backgroundTone2: "#011",
                            backgroundTone3: "#001"
                        }}
                    ></Button2>
                </div>
            {:else}
                <div in:fade>
                    <Button2 
                        btnText="Follow"
                        onClick={() => handleClick(!(user.Privacy === 1) ? 1 : 0, user)}
                        styleConfig={{ 
                            btnType: 1, 
                            btnWidth: 140,
                            btnHeight: 42,
                            fontSize: 15, 
                            fontColor: "#011",
                            fontWeight: "450",
                            hoverFontColor: "greenyellow",
                            borderRadius: 10,
                            borderWidth: 2,
                            borderColor: "greenyellow",
                            borderHoverColor: "green",
                            backgroundTone1: "greenyellow", 
                            backgroundTone2: "#011",
                            backgroundTone3: "#001"
                        }}
                    ></Button2>
                </div>
            {/if}
        </div>
    </div>
    <div class="sub-container"></div>
</div>

<style>
    .container {
        display: flex;
        justify-content: center;
        align-items: top;
        width: 100%;
        height: 290px;
    }

    .sub-container {
        display: flex;
        align-items: center;
        flex-direction: column;
        justify-content: center;
    }
</style>