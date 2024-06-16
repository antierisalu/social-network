<script>
import { createEventDispatcher } from 'svelte';
    import { userProfileData } from "../../stores"
    import { slide, fade } from 'svelte/transition'

async function selectUser(userID){
        const response = await fetch('http://localhost:8080/user?id='+userID);
        if (response.ok) {
            const selectedUser = await response.json()
            console.log("OUUU",selectedUser)
            userProfileData.set(selectedUser)
        } else {
            console.error('Error fetching users:', response.status);
        }
    }
    const dispatch = createEventDispatcher();
    export let followers
    function closeOverlay() {
      dispatch('close');
    }
</script>


<div out:fade={{duration:150}} class= "overlay" on:click={closeOverlay}></div>
<div in:slide out:fade={{duration:150}} class = "modal">
    <div class="modal-content">

        {#each followers as follower}
        <!-- svelte-ignore a11y-click-events-have-key-events -->
        <div class="singleUser" on:click={() => selectUser(follower.ID)}>
            <!-- svelte-ignore a11y-missing-attribute -->
            <img src={follower.Avatar}/> {follower.FirstName} {follower.LastName}
        </div>
        {/each}
    </div>
    <button on:click={closeOverlay}>Close Overlay</button>
</div>

<style>
img {
    max-height: 20px;
}

.overlay {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(0, 0, 0, 0.5);
    z-index: 0;
}

.singleUser {
        border: solid 1px #333;
        border-radius: 6px;
        cursor: pointer;
        padding: 8px;
        background: #011;
        margin: 2px;
    }

.modal {
    position: fixed;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    background-color: white;
    padding: 20px;
    box-shadow: 0 5px 15px rgba(0, 0, 0, 0.3);
    z-index: 2;
    width: 80%;
    max-width: 500px;
    border-radius: 8px;
}

.modal-content {
    position: relative;
}
</style>