interface CarModelProps {
  model: string;
  orSimilarText: string;
}
export function CarModel({ model, orSimilarText }: CarModelProps) {
  return (
    <div className="flex items-center gap-2 my-4">
      <h4 className="type-h4 text-navy">{model}</h4>
      <span className="text-xl font-normal text-navy">|</span>
      <span className="type-paragraph text-navy">{orSimilarText}</span>
    </div>
  );
}
