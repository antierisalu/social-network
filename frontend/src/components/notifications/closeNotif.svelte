<div class="close-btn" >
    <button on:click
        class="close-button"
    >
        <svg
            xmlns="http://www.w3.org/2000/svg"
            x="0px"
            y="0px"
            width="100"
            height="100"
            viewBox="0 0 50 50"
        >
            <path
                d="M 7.71875 6.28125 L 6.28125 7.71875 L 23.5625 25 L 6.28125 42.28125 L 7.71875 43.71875 L 25 26.4375 L 42.28125 43.71875 L 43.71875 42.28125 L 26.4375 25 L 43.71875 7.71875 L 42.28125 6.28125 L 25 23.5625 Z"
            ></path>
        </svg>
    </button>
</div>

<style>
    button {
        border: none;
        border-radius: 4px;
        cursor: pointer;
        background-color: green;
        color: white;
    }

    button:hover {
        background-color: #ddd;
        color: black;
    }

     .close-button {
        font-size: 12px;
        background-color: transparent;
        color: black;
        margin: 0px;
    }

    .close-button svg {
        width: 18px;
        height: 18px;
    }

    .close-btn {
        margin-bottom: -16px;
        display: flex;
        justify-content: flex-end;
        height: auto;
    }

    .close-btn:hover {
        color: black;
    }
</style>