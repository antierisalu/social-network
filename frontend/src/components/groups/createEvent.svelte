<script>
  import { createEventDispatcher } from "svelte";
  import { slide, fade } from "svelte/transition";
  import { userInfo, groupSelected, API_URL } from "../../stores";
  import { getEvents, sendRSVP } from "../../utils";
  import Button from "../../shared/button.svelte";

  const dispatch = createEventDispatcher();
  function closeOverlay() {
    dispatch("close");
  }

  let today = new Date().toISOString().slice(0, 10);
  let currentTime = new Date().toLocaleTimeString([], {
    hour: "2-digit",
    minute: "2-digit",
  });
  console.log(currentTime);
  function autoResize() {
    // for automatic resize of post content textarea
    const maxHeight = window.innerHeight * 0.8;
    const minHeight = 100;
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

  function toggleRSVP() {
    if (rsvpStatus === "Not Going") {
      rsvpStatus = "Not Sure";
      event.certainty = 50
    } else if (rsvpStatus === "Not Sure") {
      rsvpStatus = "Going";
      event.certainty = 90
    } else if (rsvpStatus === "Going") {
      rsvpStatus = "Not Going";
      event.certainty = 10
    }
  }

  let description = "";
  let title = "";
  let selectedDate = today;
  let selectedTime = currentTime;
  let ownerID = $userInfo.id;
  let rsvpStatus = "Going"
  let certainty = 100

  function handleDateChange(event) {
    event.date = event.target.value;
  }

  function handleTimeChange(event) {
    event.time = event.target.value;
  }

  $: event = {
    ownerID: ownerID,
    title: title,
    description: description,
    date: selectedDate,
    time: selectedTime,
    certainty: certainty
  };

  async function sendEvent() {
    if (!event.title || !event.description || !event.date || !event.time) {
      alert("Title, Description & Date cannot be empty!");
      return;
    }

    const currentDate = new Date();
    const eventDate = new Date(event.date + "T" + event.time);
    const utcEventDateStr = eventDate
      .toISOString()
      .slice(0, 19)
      .replace("T", " ");
    if (eventDate < currentDate) {
      alert("Event cannot be in the past!");
      return;
    }

    const response = await fetch(`${API_URL}/newEvent`, {
      credentials: "include",
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        groupID: $groupSelected,
        ownerID: event.ownerID,
        title: event.title,
        description: event.description,
        date: utcEventDateStr,
      }),
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    event.id = await response.json()
    
    sendRSVP(event)
    closeOverlay();
    getEvents($groupSelected);
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
        <textarea
          on:input={autoResize}
          bind:value={event.description}
          placeholder="Event Description"
        ></textarea>
      </div>
      <div class="eventDate">
        <input
          type="date"
          min={today}
          bind:value={event.date}
          on:input={handleDateChange}
        />

        <input
          type="time"
          min={currentTime}
          bind:value={event.time}
          on:input={handleTimeChange}
        />

        <Button type="secondary" inverse on:click={() => toggleRSVP()}
          >{rsvpStatus}</Button
        >
        <!--         <p>Event on: {event.date} at {event.time}</p> -->
        <p></p>
        event will be deleted 2 hours after start time
        <p></p>
      </div>
      <div class="eventButtons">
        <Button type="secondary" on:click={() => sendEvent()}
          >Create event</Button
        >
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
