import { create } from "zustand";

interface SocketState {
  socket: WebSocket | null;
  messages: string[];
  actions: {
    connect: (url: string) => void;
    send: (message: string) => void;
  };
}

interface Packet {
  topic: string;
  payload: string;
}

export const useSocketStore = create<SocketState>()((set, get) => ({
  socket: null,
  messages: [],
  actions: {
    connect: (url) => {
      if (get().socket) {
        return;
      }

      const socket = new WebSocket(url);

      socket.onopen = () => {
        set({ socket });
        console.log("WebSocket connection opened");
      };

      socket.onmessage = (event) => {
        const packet: Packet = JSON.parse(event.data);
        console.log("packet", packet);
        const messages =
          packet.topic === "import:begin" || packet.topic === "prune:begin"
            ? [packet.payload]
            : [...get().messages, packet.payload];
        set(() => ({ messages }));
      };

      socket.onclose = () => {
        console.log("WebSocket connection closed");
        set({ socket: null });
      };
    },
    send: (message) => {
      get().socket?.send(JSON.stringify(message));
    },
    close: () => {
      get().socket?.close();
    },
  },
  // connect: (url) => {
  //   const socket = new WebSocket(url);

  //   socket.onopen = () => {
  //     set({ socket });
  //     console.log("WebSocket connection opened");
  //   };

  //   socket.onmessage = (event) => {
  //     const message = JSON.parse(event.data);
  //     set((state) => ({ messages: [...state.messages, message] }));
  //   };

  //   socket.onclose = () => {
  //     console.log("WebSocket connection closed");
  //     set({ socket: null });
  //   };
  // },
  // sendMessage: (message) => {
  //   if (useWebSocketStore.getState().socket) {
  //     useWebSocketStore.getState().socket.send(JSON.stringify(message));
  //   }
  // },
}));

export const useSocketActions = () => useSocketStore((state) => state.actions);
