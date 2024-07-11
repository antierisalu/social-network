<script>
    import Button from "../../shared/button.svelte";
    import { sendRSVP } from "../../utils";
    export let event;

    console.log(event.date);
    

    let oldGoing = event.going;


    event.certainty = event.going * 40 + 50;
    let [ng, g, ns] = ["", "", ""];

    $: {
        switch (event.going) {
            case 1:
                event.certainty = 100;
                break;
            case 0:
                event.certainty = 50;
                break;
            case -1:
                event.certainty = 0;
                break;
            default:
                event.going = null;
        }

        if (event.certainty < 32) {
            [ng, g, ns] = ["selected", "", ""];
            if (oldGoing === 1) event.goingCount--, event.notGoingCount++;
            if (oldGoing === 0) event.notSureCount--, event.notGoingCount++;
            oldGoing = -1;
        } else if (event.certainty > 77) {
            [ng, g, ns] = ["", "selected", ""];
            if (oldGoing === 0) event.notSureCount--, event.goingCount++;
            if (oldGoing === -1) event.notGoingCount--, event.goingCount++;
            oldGoing = 1;
        } else {
            [ng, g, ns] = ["", "", "selected"];
            if (oldGoing === 1) event.goingCount--, event.notSureCount++;
            if (oldGoing === -1) event.notGoingCount--, event.notSureCount++;
            oldGoing = 0;
        }
    }
</script>

<div class="singleEvent">
    <div class="eventBody">
        <div class="eventInfo">
            <div class="eventTitle">{event.title}</div>

            <div class="date">
                {new Date(event.date).toLocaleString("default", {
                    weekday: "short",
                    month: "long",
                    day: "numeric",
                    year: "numeric",
                })}
                at {new Date(event.date).toLocaleString("default", {
                    hour: "numeric",
                    minute: "numeric",
                })}
            </div>
            <div class="eventDescription">{event.description}</div>
        </div>
        <div class="eventDate"></div>
    </div>
    <div class="goingStats">
        <div class="goingStatsWRAPPER">
            <div class="notGoings {ng}">Not Going: {event.notGoingCount}</div>
            <div class="notSures {ns}">Not Sure: {event.notSureCount}</div>
            <div class="goings {g}">Going: {event.goingCount}</div>
        </div>
        {#if event.going === undefined || event.going === null}
            <div class="slaider">
                <input
                    bind:value={event.certainty}
                    type="range"
                    min="0"
                    max="100"
                    class="slider"
                    style="padding: 0;"
                />
                <p class="slideText">Slide for RSVP</p>
                <div class="submitButton">
                    <Button
                        type="secondary"
                        inverse
                        on:click={() =>
                            sendRSVP(event).then((going) => {
                                event.going = going;
                            })}>Submit</Button
                    >
                </div>
            </div>
        {:else}
            <div class="going">
                <Button
                    type="secondary"
                    inverse
                    on:click={() => (event.going = null)}>Edit</Button
                >
            </div>
        {/if}
    </div>
</div>

<style>
    .selected {
        color: greenyellow;
        font-weight: bold;
    }

    .submitButton {
        text-align: center;
    }
    .goingStats {
        display: flex;
        border-top: solid 1px #555;
        padding: 5px;
        justify-content: space-between;
        flex-direction: column;
    }
    .goingStatsWRAPPER {
        display: flex;
        padding: 5px;
        justify-content: space-between;
        flex-direction: row;
    }
    .date {
        color: rgb(122, 184, 30);
        font-size: x-small;
        padding: 8px 0 4px 0;
    }

    .slideContainer {
        font-weight: bold;
        font-size: large;
        accent-color: greenyellow;
    }

    .slideText,
    .going {
        font-weight: 500;
    }

    .goingText {
        font-weight: bold;
        color: greenyellow;
    }

    p,
    .slider {
        font-weight: bold;
        font-size: large;
        accent-color: greenyellow;
        margin: 0;
        width: 100%;
    }
    .eventBody {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin: 4px 0;
    }
    .singleEvent {
        background-color: #011;
        border: solid 1px #555;
        border-radius: 8px;
        margin: 8px 0;
    }

    .eventTitle {
        font-size: large;
        font-weight: bold;
    }
    .eventInfo {
        width: 100%;
    }
</style>
