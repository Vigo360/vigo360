/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package public

import (
	"errors"
	"math/rand"

	"vigo360.es/new/internal/database"
	"vigo360.es/new/internal/model"
)

type Sugerencia struct {
	model.Publicacion
	Points int
}

// Compares both slices for matching tags, and returns how many match
func FindMatchingTags(tags1, tags2 []string) int {
	// Make it a map to avoid passing it many times
	var tag1map = make(map[string]bool)
	for _, v := range tags1 {
		tag1map[v] = true
	}

	// Iterate tag2 to find which match
	var matches = 0

	for _, v := range tags2 {
		if _, ok := tag1map[v]; ok {
			matches++
		}
	}

	return matches
}

func generateSuggestions(original_id string) ([]Sugerencia, error) {
	var resultado = make([]Sugerencia, 3)
	var ps = model.NewPublicacionStore(database.GetDB())

	original, err := ps.ObtenerPorId(original_id, true)
	if err != nil {
		return resultado, err
	}

	options, err := ps.Listar()
	if err != nil {
		return resultado, err
	}
	options = options.FiltrarPublicas()

	var originalTags []string
	for _, t2 := range original.Tags {
		originalTags = append(originalTags, t2.Id)
	}

	var pointedOptions []Sugerencia
	for _, rp := range options {
		if rp.Id == original.Id {
			continue
		}

		var np = Sugerencia{Publicacion: rp}
		var points = 10

		if original.Autor.Id == rp.Autor.Id {
			points += 12
		}

		var tags []string
		for _, t2 := range rp.Tags {
			tags = append(tags, t2.Id)
		}
		matches := FindMatchingTags(tags, originalTags)
		if len(originalTags) == matches {
			points += len(originalTags) * 4
		}

		points += matches * 3

		points += rand.Intn(12) - 6

		np.Points = points
		pointedOptions = append(pointedOptions, np)
	}

	if len(pointedOptions) < 3 {
		return resultado, errors.New("insufficient posts to recommend")
	}

	resultado[0] = pointedOptions[0]
	for _, rp := range pointedOptions[1:] {
		if resultado[0].Points < rp.Points {
			resultado[1] = resultado[0]
			resultado[0] = rp
			continue
		}

		if resultado[1].Points < rp.Points {
			resultado[1] = rp
		}
	}

	// Sugerencia al azar
	var posicionAlAzar = func() int { return rand.Intn(len(pointedOptions)) }
	resultado[2] = pointedOptions[posicionAlAzar()]

	// Comprobar que no haya ninguna duplicada
	var ids = map[string]bool{}
	for i, r := range resultado {
		var duplicado = false
		if _, ok := ids[r.Id]; ok { // Duplicado
			duplicado = true
		}

		for duplicado {
			resultado[i] = pointedOptions[posicionAlAzar()]
			if _, ok := ids[resultado[i].Id]; !ok {
				duplicado = false
			}
		}

		ids[r.Id] = true
	}

	return resultado, nil
}
