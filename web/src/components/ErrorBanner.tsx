import { type Component, type JSX, createSignal, onMount } from "solid-js";

export interface ErrorBannerProps {
  errorSearchParam: string | null;
}

const ErrorBanner: Component<ErrorBannerProps> = (props) => {
  const [show, setShow] = createSignal(true);
  const { errorSearchParam } = props;

  onMount(() => {
    console.log("setting timeout");
    setTimeout(() => {
      console.log("executing timeout callback");
      setShow(() => false);
    }, 1_000);
  });

  if (errorSearchParam === null) {
    return null;
  }

  const template = () => show() && renderErrorTemplate(errorSearchParam);

  function renderErrorTemplate(errorSearchParam: string): JSX.Element {
    return {
      "401": (
        <div class="py-2 pl-4 fixed bottom-0 left-0 right-0 bg-yellow-600">
          <i class="fa-solid fa-triangle-exclamation mr-1 align-baseline" />
          <span class="font-medium align-bottom">Please sign in again.</span>
        </div>
      ),
      "500": (
        <div class="py-2 pl-4 fixed bottom-0 left-0 right-0 bg-red-600">
          <i class="fa-solid fa-circle-exclamation mr-1 align-baseline" />
          <span class="font-medium align-bottom">
            Server error occurred. Try again, or call support if problem
            persists.
          </span>
        </div>
      ),
    }[errorSearchParam];
  }

  return template();
};

export default ErrorBanner;
