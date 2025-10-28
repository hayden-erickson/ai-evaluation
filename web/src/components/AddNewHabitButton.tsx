type Props = { onClick: () => void }
export default function AddNewHabitButton({ onClick }: Props) {
  return <button className="button" onClick={onClick}>Add New Habit</button>
}
