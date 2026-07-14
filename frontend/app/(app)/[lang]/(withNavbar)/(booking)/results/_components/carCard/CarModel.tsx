interface CarModelProps {
  model: string;
  orSimilarText: string;
}

export function CarModel({ model, orSimilarText }: CarModelProps) {
  const isLongModel = model.length > 15;

  return (
    <div
      className={`flex ${
        isLongModel ? "flex-col items-start lg:max-w-3/4" : "items-center lg:my-4 gap-2"
      }`}
    >
      <h5 className="type-h5 text-navy">{model}</h5>

      {!isLongModel && (
        <span className="text-xl font-normal text-navy">|</span>
      )}

      <span className="type-paragraph text-navy">{orSimilarText}</span>
    </div>
  );
}