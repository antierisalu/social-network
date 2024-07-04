<script>
    import { createEventDispatcher } from "svelte";
    import { slide, fade } from "svelte/transition";
    import { userInfo } from "../../stores";
    import Button from "../../shared/button.svelte";
  
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
  
    let description = "";
    let title = "";
    let selectedDate = ''; 

    function handleDateChange(event) {
    selectedDate = event.target.value;
    }
  
    $: event = {
      ID: 0,
      ownerID: $userInfo.id,
      title: title,
      description: description,
      date: selectedDate,
    };

    async function sendEvent() {
      console.log("i want to create", event);
      if (!event.title || !event.description || !event.date) {
        alert("Title, Description & Date cannot be empty!");
        return;
      }
      const response = await fetch("/newEvent", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          ID: eventID,
          ownerID: event.ownerID,
          Title: event.title,
          Description: event.description,
          date: selectedDate,
        }),
      });
  
      const eventID = await response.json();
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      closeOverlay();
      getEvents();
    }
  
  </script>
  
  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <div out:fade={{ duration: 150 }} class="overlay"></div>
  <div in:slide out:fade={{ duration: 150 }} class="modal">
    <div class="modal-content">
      <div class="createEvent">
        <div class="eventTitle">
            <textarea bind:value={event.title} placeholder="Event Title"></textarea>
        </div>
        <div class="eventDescription">
            <textarea on:input={autoResize} bind:value={event.description} placeholder="Event Description"></textarea>
        </div>
        <div class="eventDate">
            <input type="date" bind:value={selectedDate} on:input={handleDateChange} />
            <p>Event on: {selectedDate}</p>
        </div>
        <div class="eventButtons">
          <Button type="secondary" on:click={() => sendEvent()}>Create event</Button>
          <Button on:click={closeOverlay}>Cancel</Button>
        </div>
      </div>
    </div>
  </div>
  
  <style>

    p {
        margin-top: 0;
    }

    .eventTitle textarea {
      width: 100%;
    }

    .eventDescription textarea {
      min-height: 100px;
      width: 100%;
    }
  
    .createEvent {
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

    input {
  color-scheme: dark;
}
  </style>
  