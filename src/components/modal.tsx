export default function Modal(
  props: {
    children: React.ReactNode;
  },
) {
  const classes = "fixed top-0 right-0 left-0 bg-slate-800/70 size-full";

  return (
    <div className={classes}>
      {props.children}
    </div>
  );
}
