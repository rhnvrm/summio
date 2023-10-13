import { toast } from "@zerodevx/svelte-toast";

export const success = m =>
  toast.push(m, {
    theme: {
      "--toastBackground": "white",
      "--toastColor": "green",
      "--toastBarBackground": "green",
    },
  });

export const warning = m =>
  toast.push(m, {
    theme: {
      "--toastBackground": "white",
      "--toastColor": "yellow",
      "--toastBarBackground": "orange",
    },
  });

export const failure = m =>
  toast.push(m, {
    theme: {
      "--toastBackground": "white",
      "--toastColor": "red",
      "--toastBarBackground": "maroon",
    },
  });

