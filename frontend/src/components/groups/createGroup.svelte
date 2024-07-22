<script>
    import { createEventDispatcher } from "svelte";
    import { slide, fade } from "svelte/transition";
    import { userInfo,  uploadImageStore, API_URL } from "../../stores";
    import Button from "../../shared/button.svelte";
    import ImageToGroup from "../../shared/imagePreview.svelte";
    import { getGroups } from "../../utils"
  
    const dispatch = createEventDispatcher();
    function closeOverlay() {
      dispatch("close");
    }
  
    function autoResize() {
      // for automatic resize of post content textarea
      const maxHeight = window.innerHeight * 0.8;
      const minHeight = 200;
      this.style.height = "auto";
      if (this.scrollHeight > maxHeight) {
        this.style.height = maxHeight + "px";
        this.style.overflowY = "scroll";
      } else if (this.scrollHeight < minHeight) {
        this.style.height = minHeight + "px";
        this.style.overflowY = "hidden";
      } else {
        this.style.height = this.scrollHeight + "px";
        this.style.overflowY = "hidden";
      }
    }
  
    // let privatePost = false;
    // let chooseUsers = false;
    let selectedUserIds;
    let description = "";
    let title = "";
  
    let uploadImage;
    uploadImageStore.subscribe((value) => {
      uploadImage = value;
    });
  
    $: group = {
      ownerID: $userInfo.id,
      title: title,
      description: description,
      img: "",
      groupID: 0,
    //   customPrivacyIDs: selectedUserIds,
    };
  
    async function sendGroup() {
      if (!group.title || !group.description) {
        alert("Content or title cannot be empty");
        return;
      }
      const response = await fetch(`${API_URL}/newGroup`, {
        credentials: "include",
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          ownerID: group.ownerID,
          Title: group.title,
          Description: group.description,
          Img: group.img,
          GroupID: group.ID,
        }),
      });
  
      const responseGroup = await response.json();
      console.log(responseGroup)
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      closeOverlay();
      uploadImage({ group: responseGroup.id }).catch((error) => {
        console.error("Error uploading the image:", error);
      });
      getGroups();
    }
  
    function toggleUsersList() {
      chooseUsers = !chooseUsers;
    }
  </script>
  
  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <div out:fade={{ duration: 150 }} class="overlay"></div>
  <div in:slide out:fade={{ duration: 150 }} class="modal">
    <div class="modal-content">
      <div class="createGroup">
        <div class="groupTitle">
            <textarea bind:value={title} placeholder="Group Title"></textarea>
        </div>
        <div class="groupDescription">
            <textarea on:input={autoResize} bind:value={description} placeholder="Group Description"></textarea>
        </div>
        <ImageToGroup inputIDProp="groupImage" fakeInputText="Add Image" />
        <div class="groupButtons">
          <Button type="secondary" on:click={() => sendGroup()}>Create Group</Button>
          <Button on:click={closeOverlay}>Cancel</Button>
        </div>
      </div>
    </div>
  </div>
  
  <style>

    .groupTitle textarea {
      width: 100%;
    }

    .groupDescription textarea {
      min-height: 200px;
      width: 100%;
    }
  
    .createGroup {
      display: flex;
      flex-direction: column;
      border-radius: 16px;
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
  
    .modal {
      position: fixed;
      top: 50%;
      left: 50%;
      transform: translate(-50%, -50%);
      background-color: #011;
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
  