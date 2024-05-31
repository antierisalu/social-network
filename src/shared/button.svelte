<script>
    export let type = 'primary'
    export let flat = false
    export let inverse = false
    export let customStyle = ''

</script>

<button class={type} class:flat={flat} class:inverse={inverse} style={customStyle} on:click>
    <slot></slot>
</button>

<style>
    button {
        border: 0;
        cursor: pointer;
        border-radius: 6px;
        padding: 8px 12px;
        font-weight: bold;
        box-shadow: 1px 2px 3px rgba(0, 0, 0, 0.2);
    }
    .btn200px {
        width: 200px;
    }

    .primary {
        background:crimson;
        color: white;   
    }
    .secondary {
        background:greenyellow ;
        color:darkslategrey
    }
    .flat {
        box-shadow: none;
    }

    .primary.inverse {
        color: crimson;
        background: white;
        border: 2px solid crimson;
    }

    .secondary.inverse {
        color: greenyellow;
        background-color: darkslategrey;
        border: 2px solid greenyellow
    }
    
</style>